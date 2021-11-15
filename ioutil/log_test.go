package ioutil

import (
	"testing"
)

func TestFatal(t *testing.T) {
	fatal(func(i int) {}, "log: %v\n", "hello world")
}
