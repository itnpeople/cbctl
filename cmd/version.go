package cmd

import (
	"fmt"

	"github.com/itnpeople/cbctl/app"
	"github.com/spf13/cobra"
)

var (
	example = `  cbctl version`
)

// a struct to support command
type VersionOptions struct {
	app.IOStreams
}

// returns initialized Options
func NewVersionOptions(ioStreams app.IOStreams) *VersionOptions {
	return &VersionOptions{
		IOStreams: ioStreams,
	}
}

// returns a cobra command
func NewCmdVersion(ioStreams app.IOStreams) *cobra.Command {
	o := NewVersionOptions(ioStreams)
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of cbctl",
		Long:    "Print the version number of cbctl",
		Example: example,
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(o.Complete(cmd))
			app.ValidateError(o.Validate())
			app.ValidateError(o.Run())
		},
	}
	return cmd
}

// completes all the required options
func (o *VersionOptions) Complete(cmd *cobra.Command) error {
	return nil
}

// validates the provided options
func (o *VersionOptions) Validate() error {
	return nil
}

// executes version command
func (o *VersionOptions) Run() error {
	fmt.Fprintln(o.Out, fmt.Sprintf("BuildVersion=%s, BuildTime=%s", app.BuildVersion, app.BuildTime))
	return nil
}
