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
		t.Log("测试成功\n")
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
		t.Log("测试成功\n")
	})
}

func Test4(t *testing.T) {
	t.Run("test4 All功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if fmt.Sprintf("%#v", All(al)) != "[]int{1, 2, 3}" {
			t.Fatal("值错误")
		}
		t.Logf("成功：%#v\n", All(al))
	})
}

func Test5(t *testing.T) {
	t.Run("test5 GetByIndexes功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if fmt.Sprintf("%#v", GetByIndexes(al, 0, 2)) != "[]int{1, 3}" {
			t.Fatal("值错误")
		}
		t.Logf("成功：%#v\n", GetByIndexes(al, 0, 2))
	})
}
