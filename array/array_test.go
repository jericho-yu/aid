package array

import (
	"testing"
)

type A struct {
	Name string
}

func Test1(t *testing.T) {
	t.Run("Test1", func(t *testing.T) {
		t.Logf("%#v", ToAny([]A{{Name: "1"}, {Name: "2"}}))
	})
}
