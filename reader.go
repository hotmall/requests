package requests

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInsufCap = errors.New("insufficient capacity")
)

type IntReader struct {
	rdr *strings.Reader
}

func NewIntReader(v int) *IntReader {
	s := strconv.Itoa(v)
	return &IntReader{
		rdr: strings.NewReader(s),
	}
}

func (r IntReader) Read(p []byte) (n int, err error) {
	if cap(p) < int(r.rdr.Size()) {
		err = ErrInsufCap
		return
	}
	n, err = r.rdr.Read(p)
	return
}
