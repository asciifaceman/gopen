package nmap

import (
	"fmt"

	"github.com/asciifaceman/gopen/pkg/parse"
)

// Options for task
type Options struct {
	Target  string
	TmpFile string
	Flags   []string
	Content *parse.NmapRun
}

// Option functional API for setting options
type Option func(*Options)

// WithTarget adds a target
func WithTarget(target string) Option {
	return func(opt *Options) {
		opt.Target = target
	}
}

// WithFlag adds nmap flags
func WithFlag(flag string) Option {
	return func(opt *Options) {
		opt.Flags = append(opt.Flags, flag)
	}
}

// WithFlags adds a slice of flags
func WithFlags(flags []string) Option {
	return func(opt *Options) {
		var formattedFlags []string
		for _, flag := range flags {
			formattedFlags = append(formattedFlags, fmt.Sprintf("-%s", flag))
		}
		opt.Flags = formattedFlags
	}
}

// StripQuotes ...
/*
func StripQuotes(flag string) string {

	flag = strings.Replace(flag, '\"', "", -1)
	return flag
}
*/
