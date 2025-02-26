package anyMap

import "testing"

func Test1(t *testing.T) {
	t.Run("test1 初始化一个空的AnyMap", func(t *testing.T) {
		am := MakeAnyMap[int, string]()
		if am == nil {
			t.Fatal("am is nil")
		}
	})
}
