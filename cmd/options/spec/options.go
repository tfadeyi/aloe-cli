package spec

import (
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tfadeyi/aloe-cli/internal/generate"
	errhandler "github.com/tfadeyi/go-aloe"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		StdOutput    bool
		Formats      []string
		IncludedDirs []string
	}
)

// New creates a new instance of the application's options
func New() *Options {
	return new(Options)
}

// Prepare assigns the applications flag/options to the cobra cli
func (o *Options) Prepare(cmd *cobra.Command) *Options {
	o.addAppFlags(cmd.Flags())
	return o
}

// Complete initialises the components needed for the application to function given the options
func (o *Options) Complete() error {
	for i, format := range o.Formats {
		cFormat := strings.ToLower(strings.TrimSpace(format))
		if !generate.IsValidOutputFormat(cFormat) {
			// @aloe code invalid_output_format
			// @aloe title invalid_output_format
			// @aloe summary the output format passed to --format was invalid, valid: json,yaml
			err := errors.Errorf("the output format given %q is not valid", format)
			return errhandler.DefaultOrDie().Error(err, "invalid_output_format")
		}

		o.Formats[i] = cFormat
	}

	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(
		&o.IncludedDirs,
		"dirs",
		[]string{"./."},
		"Comma separated list of directories to be parses by the tool",
	)
	fs.StringSliceVar(
		&o.Formats,
		"format",
		[]string{"yaml"},
		"Output format (yaml,json,markdown)",
	)
	fs.BoolVar(
		&o.StdOutput,
		"stdout",
		false,
		"Print output to standard output.",
	)
}
