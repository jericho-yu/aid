package steam

import (
	"bytes"
	"errors"
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
func (my *Steam) Copy(fn func(copied []byte) error) (io.ReadCloser, error) {
	var (
		err    error
		copied []byte
	)

	if my.readCloser == nil {
		return nil, errors.New("空内容")
	}

	copied, err = io.ReadAll(my.readCloser)
	if err != nil {
		return nil, err
	}

	err = fn(copied)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewBuffer(copied)), err
}
