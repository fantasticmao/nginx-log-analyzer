package ioutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenFile(t *testing.T) {
	file, isGzip := OpenFile("../testdata/access.log")
	assert.NotNil(t, file)
	assert.False(t, isGzip)

	file, isGzip = OpenFile("../testdata/access.json.log")
	assert.NotNil(t, file)
	assert.False(t, isGzip)

	file, isGzip = OpenFile("../testdata/access.json.log.1.gz")
	assert.NotNil(t, file)
	assert.True(t, isGzip)
}

func TestReadFile(t *testing.T) {
	file, isGzip := OpenFile("../testdata/access.log")
	reader := ReadFile(file, isGzip)
	assert.NotNil(t, reader)

	file, isGzip = OpenFile("../testdata/access.json.log")
	reader = ReadFile(file, isGzip)
	assert.NotNil(t, reader)

	file, isGzip = OpenFile("../testdata/access.json.log.1.gz")
	reader = ReadFile(file, isGzip)
	assert.NotNil(t, reader)
}
