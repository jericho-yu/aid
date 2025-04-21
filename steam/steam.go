package steam

import "io"

type (
	Steam struct {
		readCloser io.ReadCloser
		err        error
	}
)

var (
	App Steam
)

// New 实例化：Steam
func (*Steam) New(readCloser io.ReadCloser) *Steam { return &Steam{readCloser: readCloser} }

// Copy 复制流
func (my *Steam) Copy(fn func(copied []byte) error) *Steam {
	var b []byte
	b, my.err = io.ReadAll(my.readCloser)
	if my.err != nil {
		return my
	}

	my.err = fn(b)

	return my
}
