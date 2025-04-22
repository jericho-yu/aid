package steam

import (
	"bytes"
	"io"
)

type (
	Steam struct {
		readCloser io.ReadCloser
		Error      error
	}
)

var (
	App Steam
)

// New 实例化：Steam
func (*Steam) New(readCloser io.ReadCloser) *Steam { return &Steam{readCloser: readCloser} }

// Copy 复制流
func (my *Steam) Copy(fn func(copied []byte) error) *Steam {
	if my.readCloser == nil {
		return my
	}

	var b []byte
	b, my.Error = io.ReadAll(my.readCloser)
	if my.Error != nil {
		return my
	}

	my.Error = fn(b)
	if my.Error != nil {
		return my
	}

	my.readCloser = io.NopCloser(bytes.NewBuffer(b))

	return my
}
