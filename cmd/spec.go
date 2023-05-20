package cmd

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	specoptions "github.com/tfadeyi/aloe-cli/cmd/options/spec"
	"github.com/tfadeyi/aloe-cli/internal/generate"
	"github.com/tfadeyi/aloe-cli/internal/logging"
	"github.com/tfadeyi/aloe-cli/internal/parser/golang"
	errhandler "github.com/tfadeyi/go-aloe"
)

// specCmd is the entrypoint to the specification sub commands
func specCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec",
		Short: "Command to operate on the aloe specification",
		Long:  ``,
	}
	cmd.AddCommand(specGenerateCmd(), specValidateCmd())
	return cmd
}

func specGenerateCmd() *cobra.Command {
	opts := specoptions.New()
	cmd := &cobra.Command{
		Use:           "generate",
		Short:         "Generates the aloe specification from a given source code",
		Long:          ``,
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Complete()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			source := args[0]
			//output := specoptions.defaultOutputFile
			//if len(args) == 2 {
			//	// the output file path was passed
			//	output = args[1]
			//}

			// @aloe code clean_artifacts_error
			// @aloe title Error Removing Previous Artifacts
			// @aloe summary The tool has failed to delete the artifacts from the previous execution.
			// @aloe details The tool has failed to delete the artifacts from the previous execution.
			// Try manually deleting them before running the tool again.

			goparser := golang.New(logger, source, opts.IncludedDirs...)
			app, err := goparser.Parse()
			if err != nil {
				return err
			}
			//wr, err := outputWriter(opts.StdOutput)
			//if err != nil {
			//	return err
			//}

			return generate.WriteSpecification(app, opts.StdOutput, opts.Formats...)
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func specValidateCmd() *cobra.Command {
	opts := specoptions.New()
	cmd := &cobra.Command{
		Use:           "validate",
		Short:         "Validates a given aloe specification",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// @aloe code validate_not_implemented
			// @aloe title validate_not_implemented
			// @aloe summary spec validate command has not been implemented yet
			// @aloe details specification validate command has not been implemented yet, will be implemented shortly
			return errhandler.DefaultOrDie().Error(errors.New("not implemented"), "validate_not_implemented")
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand(specCmd())
}
