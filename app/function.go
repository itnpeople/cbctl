package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DefaultErrorExitCode = 1
)

func DefaultSubCommandRun(out io.Writer) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		c.SetOutput(out)
		RequireNoArguments(c, args)
		c.Help()
	}
}

// RequireNoArguments exits with a usage error if extra arguments are provided.
func RequireNoArguments(c *cobra.Command, args []string) {
	if len(args) > 0 {
		ValidateError(UsageErrorf(c, "unknown command %q", strings.Join(args, " ")))
	}
}

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func ValidateError(err error) {

	if err != nil {
		msg := err.Error()
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
		os.Exit(DefaultErrorExitCode)
	}

}
