package nmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/asciifaceman/gopen/pkg/parse"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	execCmd  = exec.Command
	nmapPath = "/usr/local/bin/nmap"
)

// Config ...
type Config struct {
	Logger   *zap.Logger
	Size     int
	ExecPath string
}

// Scanner ...
type Scanner struct {
	Logger   *zap.Logger
	ExecPool sync.WaitGroup
	Pending  chan *Task
	//Complete    chan *Task
	fs          afero.Fs
	outputFlags []string
}

// GetFormattedFlags ...
func (s *Scanner) GetFormattedFlags(t *Task) ([]string, error) {
	var response []string
	for _, flag := range t.Flags {
		if strings.Contains(flag, " ") {
			flg := strings.Split(flag, " ")
			if len(flg) > 1 {
				s.Logger.Error("Too long of a flag argument",
					zap.String("flag", flag),
				)
				return []string{}, errors.New("too long of a flag argument")
			}
			for _, f := range flg {
				response = append(response, f)
			}
		} else {
			response = append(response, flag)
		}
	}

	for _, flag := range s.outputFlags {
		response = append(response, flag)
	}
	response = append(response, t.Target)
	return response, nil
}

// NewScanner ...
func NewScanner(c *Config) (*Scanner, error) {
	pending := make(chan *Task, 32)

	// Create object
	s := &Scanner{
		Logger:  c.Logger,
		Pending: pending,
		outputFlags: []string{
			"-oX", // output xml
			"-",   // to stdout
		},
		fs: afero.NewMemMapFs(),
	}

	// Launch pool
	for w := 1; w <= c.Size; w++ {
		s.ExecPool.Add(1)
		go s.ScanPool(w)
	}

	return s, nil
}

// Task ...
type Task struct {
	Target  string
	TmpFile string
	Flags   []string
	Content *parse.NmapRun
}

// NewTaskWithOptions ...
func (s *Scanner) NewTaskWithOptions(optFnc ...Option) error {
	if len(optFnc) < 1 {
		return errors.New("nothing to set up")
	}

	options := &Options{}
	for _, opt := range optFnc {
		opt(options)
	}

	t := &Task{
		Target: options.Target,
		Flags:  options.Flags,
	}
	s.Pending <- t
	return nil
}

// NewTaskWithDefaultOptions ...
func (s *Scanner) NewTaskWithDefaultOptions(target string) error {
	t := &Task{
		Target: target,
		Flags: []string{
			"-v",
			"-A",
		},
	}
	s.Pending <- t
	return nil
}

// ScanPool ...
func (s *Scanner) ScanPool(id int) error {
	defer s.ExecPool.Done()
	s.Logger.Info("has joined the game",
		zap.Int("scanner", id),
	)

	//pending, ok <-s.Pending

	for j := range s.Pending {
		// Time lock
		start := time.Now()

		s.Logger.Info("received task",
			zap.Int("scanner", id),
			zap.String("target", j.Target),
		)

		// Scan host and output to temp file
		formatted, err := s.GetFormattedFlags(j)
		if err != nil {
			return err
		}
		cmd := execCmd(nmapPath, formatted...)
		//spew.Dump(cmd.Args)
		out, err := cmd.CombinedOutput()
		if err != nil {
			args := strings.Join(cmd.Args[:], ",")
			s.Logger.Error("failed to run namp",
				zap.Int("scanner", id),
				zap.String("target", j.Target),
				zap.String("error", err.Error()),
				zap.String("command", args),
			)
			//spew.Dump(string(out))
			return err
		}

		if len(out) < 1 {
			s.Logger.Error("nmap returned empty",
				zap.Int("scanner", id),
				zap.String("target", j.Target),
				zap.String("error", "unknown"),
			)
		}

		// Parse
		parsed, err := parse.Parse(out)
		if err != nil {
			s.Logger.Error("failed to parse namp output",
				zap.Int("scanner", id),
				zap.String("target", j.Target),
				zap.String("error", err.Error()),
			)
			return err
		}

		filename := fmt.Sprintf("%s-%v.json", j.Target, start.Second())
		s.Logger.Info("writing results to file",
			zap.String("target", j.Target),
			zap.String("filename", filename),
		)
		file, err := json.MarshalIndent(parsed, "", " ")
		if err != nil {
			s.Logger.Error("failed to marshal namp json",
				zap.Int("scanner", id),
				zap.String("target", j.Target),
				zap.String("error", err.Error()),
			)
			return err
		}

		err = ioutil.WriteFile(filename, file, 0644)
		if err != nil {
			s.Logger.Error("failed to write json results",
				zap.Int("scanner", id),
				zap.String("target", j.Target),
				zap.String("error", err.Error()),
			)
			return err
		}
		err = os.Chown(filename, os.Getuid(), os.Getgid())
		if err != nil {
			s.Logger.Error("failed to write permissions on json file",
				zap.Int("scanner", id),
				zap.String("target", j.Target),
				zap.String("error", err.Error()),
			)
			return err
		}

		// Collapse timelock and derive duration
		t := time.Now()
		duration := t.Sub(start)
		s.Logger.Info("completed task",
			zap.Int("scanner", id),
			zap.String("target", j.Target),
			zap.Duration("elapsed", duration),
		)

	}
	s.Logger.Info("has served its time",
		zap.Int("scanner", id),
	)
	return nil
}
