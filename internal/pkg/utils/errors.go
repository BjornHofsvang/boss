package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func fatal(msg string, code int) {

	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

var fatalErrHandler = fatal

const (
	DefaultErrorExitCode = 1
)

func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func CheckError(err error) {
	checkError(err, fatalErrHandler)
}

func checkError(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}

	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = fmt.Sprintf("error: %s", msg)
	}
	handleErr(msg, DefaultErrorExitCode)

}

func HelpErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	CheckError(cmd.Help())
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s", msg)
}

func HelpError(cmd *cobra.Command, args ...interface{}) error {
	CheckError(cmd.Help())
	msg := fmt.Sprint(args...)
	return fmt.Errorf("%s", msg)
}
