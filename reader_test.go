package requests

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntReader(t *testing.T) {
	items := []struct {
		v int
		n int
	}{
		{5, 1},
		{12, 2},
		{100, 3},
		{111, 3},
		{5678, 4},
		{56789, 5},
		{567890, 6},
		{567890345, 9},
		{5678903451, 10},
	}
	assert := assert.New(t)

	for _, item := range items {
		r := NewIntReader(item.v)

		p := make([]byte, 10)
		n, err := r.Read(p)
		assert.Nil(err)
		assert.Equal(item.n, n)

		n, err = r.Read(p)
		assert.ErrorIs(err, io.EOF)
		assert.Equal(0, n)
	}
}

func TestIntReader2(t *testing.T) {
	items := []struct {
		v int
		n int
	}{
		{56789034512, 10},
		{56789034512333, 10},
	}
	assert := assert.New(t)

	for _, item := range items {
		r := NewIntReader(item.v)

		p := make([]byte, 10)
		n, err := r.Read(p)
		assert.Equal(0, n)
		assert.ErrorIs(err, ErrInsufCap)
	}
}
