package requests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiForm(t *testing.T) {
	assert := assert.New(t)

	mf := make(MultiForm)
	filename := "hotmall.jpeg"
	mf["filename"] = strings.NewReader(filename)

	r, ok := mf["filename"]
	assert.True(ok)
	rdr, ok := r.(*strings.Reader)
	assert.True(ok)

	b := make([]byte, rdr.Size())
	n, err := rdr.Read(b)
	assert.Nil(err)
	assert.Equal(len(filename), n)
	filename2 := string(b)
	assert.Equal(filename, filename2)

	rdr.Reset(filename2)

	r2, ok := mf["filename"]
	assert.True(ok)
	b2 := make([]byte, len(filename))
	n2, err2 := r2.Read(b2)
	assert.Nil(err2)
	filename3 := string(b2)
	assert.Equal(len(filename), n2)
	assert.Equal(filename, filename3)
}
