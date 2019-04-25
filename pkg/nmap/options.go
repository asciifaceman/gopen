package nmap

import (
	"fmt"
)

// Options for task
type Options struct {
	Target string
	Flags  []string
}

// Option functional API for setting options
type Option func(*Options)

// WithTarget adds a target
func WithTarget(target string) Option {
	return func(opt *Options) {
		opt.Target = target
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
