package requests

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilename(t *testing.T) {
	assert := assert.New(t)

	filename := "aaaaaa.jpeg"
	mf := make(MultiForm)
	mf[FILENAME_KEY] = strings.NewReader(filename)

	filename2 := parseFilename(mf)
	assert.Equal(filename, filename2)

	r, ok := mf[FILENAME_KEY]
	assert.True(ok)

	b := make([]byte, len(filename))
	n, err := r.Read(b)
	assert.Nil(err)
	assert.Equal(len(filename), n)
	filename3 := string(b)
	assert.Equal(filename, filename3)
	n, err = r.Read(b)
	assert.ErrorIs(err, io.EOF)
	assert.Equal(0, n)
}

func TestParseFilename2(t *testing.T) {
	assert := assert.New(t)

	filename := "aaaaaa.jpeg"
	mf := make(MultiForm)
	mf["filename2"] = strings.NewReader(filename)

	filename2 := parseFilename(mf)
	assert.Equal(defaultFilename, filename2)
}
