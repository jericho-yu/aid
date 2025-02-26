package anyList

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	t.Run("test1 初始化一个长度为3的AnyList", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if al == nil {
			t.Fatal("al is nil")
		}
		t.Logf("测试成功，长度：%d\n", Len(al))
	})
}

func Test2(t *testing.T) {
	t.Run("test2 Has功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if !Has(al, 2) {
			t.Fatal("Has error")
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 Set功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		Set(al, 2, 4)
		if !Has(al, 2) {
			t.Fatal("Set error")
		}
		t.Log("测试成功\n")
		if val, exist := Get(al, 2); !exist {
			t.Fatal("Get error")
		} else {
			if val != 4 {
				t.Fatal("测试值错误")
			}
		}
	})
}

func Test4(t *testing.T) {
	t.Run("test4 All功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if fmt.Sprintf("%#v", All(al)) != "[]int{1, 2, 3}" {
			t.Fatal("值错误")
		}
	})
}

func Test5(t *testing.T) {
	t.Run("test5 GetByIndexes功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if fmt.Sprintf("%#v", GetByIndexes(al, 0, 2)) != "[]int{1, 3}" {
			t.Fatal("值错误")
		}
	})
}

func Test6(t *testing.T) {
	t.Run("test6 Append功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		Append(al, 4)
		if fmt.Sprintf("%#v", All(al)) != "[]int{1, 2, 3, 4}" {
			t.Fatal("值错误")
		}
	})
}

func Test7(t *testing.T) {
	t.Run("test7 First功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if First(al) != 1 {
			t.Fatal("值错误")
		}
	})
}

func Test8(t *testing.T) {
	t.Run("test Last功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if Last(al) != 3 {
			t.Fatal("值错误")
		}
	})
}

func Test9(t *testing.T) {
	t.Run("test FindIndex功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if FindIndex(al, func(val int) bool { return val == 2 }) != 1 {
			t.Fatal("值错误")
		}
	})
}
