package ioutil

import (
	"bufio"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
)

func OpenFile(path string) (*os.File, bool) {
	file, err := os.Open(path)
	if err != nil {
		Fatal("open file error: %v\n", err.Error())
		return nil, false
	}

	ext := filepath.Ext(file.Name())
	return file, strings.EqualFold(".gz", ext)
}

func ReadFile(file *os.File, isGzip bool) *bufio.Reader {
	if isGzip {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			Fatal("gzip new reader error: %v\n", err.Error())
			return nil
		}
		return bufio.NewReader(gzipReader)
	} else {
		return bufio.NewReader(file)
	}
}
