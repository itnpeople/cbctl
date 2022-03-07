package app

import (
	"fmt"
	"os"
	"strings"
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
