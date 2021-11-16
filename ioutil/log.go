package ioutil

import (
	"fmt"
	"io"
	"os"
)

func Fatal(format string, a ...interface{}) {
	fatal(os.Stderr, os.Exit, format, a...)
}

func fatal(w io.Writer, exit func(int), format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format, a...)
	exit(1)
}
