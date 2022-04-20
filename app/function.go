package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func ValidateError(err error) {

	if err != nil {
		msg := err.Error()
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}

}

func ValidCommandArgs(idx int, value *string) func(c *cobra.Command, args []string) error {

	return func(c *cobra.Command, args []string) error {
		if len(args) > idx && len(args[idx]) > 0 {
			*value = args[idx]
		}
		if *value == "" {
			return fmt.Errorf("arguemnt[%d] is empty", idx)
		}
		return nil
	}

}
