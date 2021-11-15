package ioutil

import (
	"fmt"
	"os"
)

func Fatal(format string, a ...interface{}) {
	fatal(os.Exit, format, a...)
}

func fatal(exit func(int), format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	exit(1)
}
