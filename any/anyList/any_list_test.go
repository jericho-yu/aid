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
		if fmt.Sprintf("%#v", ToSlice(al)) != "[]int{1, 2, 3}" {
			t.Fatal("错误")
		}
	})
}

func Test5(t *testing.T) {
	t.Run("test5 GetByIndexes功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if fmt.Sprintf("%#v", GetByIndexes(al, 0, 2)) != "[]int{1, 3}" {
			t.Fatal("错误")
		}
	})
}

func Test6(t *testing.T) {
	t.Run("test6 Append功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		Append(al, 4)
		if fmt.Sprintf("%#v", ToSlice(al)) != "[]int{1, 2, 3, 4}" {
			t.Fatal("错误")
		}
	})
}

func Test7(t *testing.T) {
	t.Run("test7 First功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if First(al) != 1 {
			t.Fatal("错误")
		}
	})
}

func Test8(t *testing.T) {
	t.Run("test Last功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if Last(al) != 3 {
			t.Fatal("错误")
		}
	})
}

func Test9(t *testing.T) {
	t.Run("test FindIndex功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if FindIndex(al, func(val int) bool { return val == 2 }) != 1 {
			t.Fatal("错误")
		}
	})
}

func Test10(t *testing.T) {
	t.Run("test FindIndexes功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if len(FindIndexes(al, func(val int) bool { return val > 1 })) != 2 {
			t.Fatal("错误")
		}
	})
}

func Test11(t *testing.T) {
	t.Run("test11 Copy功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		al2 := Copy(al)
		if (&al == &al2) && (fmt.Sprintf("%#v", ToSlice(al)) != fmt.Sprintf("%#v", ToSlice(al2))) {
			t.Fatal("错误")
		}
	})
}

func Test12(t *testing.T) {
	t.Run("test12 Shuffle功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
		Shuffle(al)
		t.Logf("%#v\n", ToSlice(al))
	})
}

func Test13(t *testing.T) {
	t.Run("test13 Len功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		if Len(al) != 3 {
			t.Fatal("错误")
		}
	})
}

func Test14(t *testing.T) {
	t.Run("test14 Filter功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3})
		al2 := Filter(al, func(val int) bool { return val > 1 })
		if Len(al2) != 2 {
			t.Fatalf("错误")
		}
	})
}

func Test15(t *testing.T) {
	t.Run("test15 RemoveEmpty功能", func(t *testing.T) {
		al := NewAnyList([]string{"1", "", "3"})
		RemoveEmpty(al)
		if Len(al) != 2 && fmt.Sprintf("%#v", ToSlice(al)) != `[]string{"1", "3"}` {
			t.Fatalf("错误")
		}
	})
}

func Test16(t *testing.T) {
	t.Run("test16 RemoveEmpty功能2", func(t *testing.T) {
		al := NewAnyList([]int{1, 0, 3})
		RemoveEmpty(al)
		if Len(al) != 2 && fmt.Sprintf("%#v", ToSlice(al)) != `[]int{1, 3}` {
			t.Fatalf("错误")
		}
	})
}

func Test17(t *testing.T) {
	t.Run("test17 RemoveEmpty功能3", func(t *testing.T) {
		type S struct{ s string }
		s1 := &S{s: "a"}
		s2 := &S{s: "c"}
		al := NewAnyList([]*S{s1, nil, s2})
		t.Logf("%#v", ToSlice(al))
		RemoveEmpty(al)
		t.Logf("%#v", ToSlice(al))
	})
}

func Test18(t *testing.T) {
	t.Run("test18 Join功能", func(t *testing.T) {
		al := NewAnyList([]string{"a", "b", "c"})
		if Join(al, ",") != "a,b,c" {
			t.Fatalf("错误")
		}
	})
}

func Test19(t *testing.T) {
	t.Run("test19 Join功能2", func(t *testing.T) {
		type S struct{ s string }
		s1 := &S{s: "a"}
		s2 := &S{s: "c"}
		al := NewAnyList([]*S{s1, nil, s2})
		t.Logf("%#v", Join(al, ","))
	})
}

func Test20(t *testing.T) {
	t.Run("test20 JoinWithoutEmpty功能", func(t *testing.T) {
		type S struct{ s string }
		s1 := &S{s: "a"}
		s2 := &S{s: "c"}
		al := NewAnyList([]*S{s1, nil, s2})
		t.Logf("%#v", JoinWithoutEmpty(al, ","))
	})
}

func Test21(t *testing.T) {
	t.Run("test21 In功能", func(t *testing.T) {
		al := NewAnyList([]string{"a", "", "c"})
		if !In(al, "a") {
			t.Logf("错误")
		}
	})
}

func Test22(t *testing.T) {
	t.Run("test22 NotIn功能", func(t *testing.T) {
		al := NewAnyList([]string{"a", "", "c"})
		if !NotIn(al, "b") {
			t.Log("错误")
		}
	})
}

func Test23(t *testing.T) {
	t.Run("test23 AllEmpty功能", func(t *testing.T) {
		al := NewAnyList([]string{"a", "", "c"})
		if AllEmpty(al) {
			t.Log("错误")
		}
	})
}

func Test24(t *testing.T) {
	t.Run("test24 AllEmpty功能2", func(t *testing.T) {
		al := NewAnyList([]string{"", "", ""})
		if !AllEmpty(al) {
			t.Log("错误")
		}
	})
}

func Test25(t *testing.T) {
	t.Run("test25 AllNotEmpty功能", func(t *testing.T) {
		al := NewAnyList([]string{"a", "", "c"})
		if !AllNotEmpty(al) {
			t.Log("错误")
		}
	})
}

func Test26(t *testing.T) {
	t.Run("test26 Chunk功能", func(t *testing.T) {
		al := MakeAnyList[int](0)
		for i := 0; i < 100; i++ {
			Append(al, i+1)
		}

		chunk := Chunk(al, 30)
		if len(chunk) != 4 {
			t.Fatal("错误")
		}
		t.Logf("%#v\n", chunk)
	})
}

func Test27(t *testing.T) {
	t.Run("test27 Pluck功能", func(t *testing.T) {
		type S struct {
			Name string
			Age  int
		}

		al := NewAnyList([]S{
			{Name: "张三", Age: 18},
			{Name: "李四", Age: 19},
			{Name: "王五", Age: 20},
			{Name: "赵六", Age: 21},
		})

		ages := Pluck[S, int](al, func(item S) int { return item.Age })

		if fmt.Sprintf("%#v", ages) != "[]int{18, 19, 20, 21}" {
			t.Fatal("错误")
		}
	})
}

func Test28(t *testing.T) {
	t.Run("test28 Unique功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 1, 2, 3})
		Unique(al)

		if fmt.Sprintf("%#v", ToSlice(al)) != "[]int{1, 2, 3}" {
			t.Fatal("错误")
		}
	})
}

func Test29(t *testing.T) {
	t.Run("test29 Unique功能2", func(t *testing.T) {
		al := NewAnyList([]string{"a", "a", "b", "c"})
		Unique(al)

		if fmt.Sprintf("%#v", ToSlice(al)) != `[]string{"a", "b", "c"}` {
			t.Fatal("错误")
		}
	})
}

func Test30(t *testing.T) {
	t.Run("test30 Unique功能3", func(t *testing.T) {
		type S struct{ s string }
		al := NewAnyList([]S{{s: "a"}, {s: "a"}, {s: "b"}, {s: "c"}})
		Unique(al)

		if fmt.Sprintf("%#v", ToSlice(al)) != `[]anyList.S{anyList.S{s:"a"}, anyList.S{s:"b"}, anyList.S{s:"c"}}` {
			t.Fatal("错误")
		}
	})
}

func Test31(t *testing.T) {
	t.Run("test31 RemoveByIndexes功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3, 4, 5})
		RemoveByIndexes(al, 1, 3)

		t.Log(fmt.Sprintf("%#v", ToSlice(al)))
		if fmt.Sprintf("%#v", ToSlice(al)) != "[]int{1, 3, 5}" {
			t.Fatal("错误")
		}
	})
}

func Test32(t *testing.T) {
	t.Run("test32 Every功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3, 4, 5})
		Every(al, func(idx int, item int) int { return item * item })

		if fmt.Sprintf("%#v", ToSlice(al)) != "[]int{1, 4, 9, 16, 25}" {
			t.Fatal("错误")
		}
	})

}

func Test33(t *testing.T) {
	t.Run("test33 Each功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3, 4, 5})
		al2 := MakeAnyList[int](5)

		Each(al, func(idx int, item int) { Set(al2, idx, item*5) })

		if fmt.Sprintf("%#v", ToSlice(al2)) != "[]int{5, 10, 15, 20, 25}" {
			t.Fatal("错误")
		}
	})
}

func Test34(t *testing.T) {
	t.Run("test34 Clean功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3, 4, 5})
		Clean(al)

		if fmt.Sprintf("%#v", ToSlice(al)) != "[]int{}" {
			t.Fatal("错误")
		}
	})
}

func Test35(t *testing.T) {
	t.Run("test35 Cast功能", func(t *testing.T) {
		al := NewAnyList([]int{1, 2, 3, 4, 5})
		al2 := Cast[int, string](al, func(value int) string { return fmt.Sprintf("%v", value) })

		if fmt.Sprintf("%#v", ToSlice(al2)) != `[]string{"1", "2", "3", "4", "5"}` {
			t.Fatal("错误")
		}
	})
}
