package anyMap

import (
	"fmt"
	"reflect"
	"testing"
)

func Test1(t *testing.T) {
	t.Run("test1 初始化一个空的AnyMap", func(t *testing.T) {
		am := MakeAnyMap[int, string]()
		if am == nil {
			t.Fatalf("错误")
		}
	})
}

func Test2(t *testing.T) {
	t.Run("test2 初始化AnyMap", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"分数":100, "年龄":18}` {
			t.Fatal("错误")
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 GetKeyByIndex功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if GetKeyByIndex(am, 1) != "分数" {
			t.Fatal("错误")
		}
	})
}

func Test4(t *testing.T) {
	t.Run("test4 GetKeysByIndexes功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", GetKeysByIndexes(am, 0, 1)) != `[]string{"年龄", "分数"}` {
			t.Fatal("错误")
		}
	})
}

func Test5(t *testing.T) {
	t.Run("test5 GetKeyByValue功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if GetKeyByValue(am, 18) != "年龄" {
			t.Fatal("错误")
		}
	})
}

func Test6(t *testing.T) {
	t.Run("test6 GetKeysByValues功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", GetKeysByValues(am, 18, 100)) != `[]string{"年龄", "分数"}` {
			t.Fatal("错误")
		}
	})
}

func Test7(t *testing.T) {
	t.Run("test7 GetValueByIndex功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if GetValueByIndex(am, 0) != 18 {
			t.Fatal("错误")
		}
	})
}

func Test8(t *testing.T) {
	t.Run("test8 GetValuesByIndexes功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", GetValuesByIndexes(am, 0, 1)) != "[]int{18, 100}" {
			t.Fatal("错误")
		}
	})
}

func Test9(t *testing.T) {
	t.Run("test9 GetValueByKey功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if GetValueByKey(am, "年龄") != 18 {
			t.Fatal("错误")
		}
	})
}

func Test10(t *testing.T) {
	t.Run("test10 GetValueByKey功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", GetValuesByKeys(am, "年龄", "分数")) != `[]int{18, 100}` {
			t.Fatal("错误")
		}
	})
}

func Test11(t *testing.T) {
	t.Run("test11 GetIndexByKey功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if GetIndexByKey(am, "分数") != 1 {
			t.Fatal("错误")
		}
	})
}

func Test12(t *testing.T) {
	t.Run("test12 GetIndexesByKeys功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", GetIndexesByKeys(am, "分数", "年龄")) != `[]int{1, 0}` {
			t.Fatal("错误")
		}
	})
}

func Test13(t *testing.T) {
	t.Run("test13 GetIndexByValue功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if GetIndexByValue(am, 100) != 1 {
			t.Fatal("错误")
		}
	})
}

func Test14(t *testing.T) {
	t.Run("test14 GetIndexesByValues功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if fmt.Sprintf("%#v", GetIndexesByValues(am, 100, 18)) != `[]int{1, 0}` {
			t.Fatal("错误")
		}
	})
}

func Test15(t *testing.T) {
	t.Run("test15 IsEmpty功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{})

		if !IsEmpty(am) {
			t.Fatal("错误")
		}
	})
}

func Test16(t *testing.T) {
	t.Run("test16 IsNotEmpty功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{})

		if IsNotEmpty(am) {
			t.Fatal("错误")
		}
	})
}

func Test17(t *testing.T) {
	t.Run("test17 Set功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"分数": 100})

		Set(am, "年龄", 18)
		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"分数":100, "年龄":18}` {
			t.Fatal("错误")
		}
	})
}

func Test18(t *testing.T) {
	t.Run("test18 Copy功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"分数": 100})
		am2 := Copy(am)

		if !reflect.DeepEqual(am, am2) {
			t.Fatal("错误")
		}
	})
}

func Test19(t *testing.T) {
	t.Run("test19 Len功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if Len(am) != 2 {
			t.Fatal("错误")
		}
	})
}

func Test20(t *testing.T) {
	t.Run("test20 First功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if k, v := First(am); k != "年龄" || v != 18 {
			t.Fatal("错误")
		}

	})
}

func Test21(t *testing.T) {
	t.Run("test21 Last功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		if k, v := Last(am); k != "分数" || v != 100 {
			t.Fatal("错误")
		}
	})
}

func Test22(t *testing.T) {
	t.Run("test22 Filter功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})

		am2 := Filter(am, func(key string, value int) bool { return value > 18 })

		if fmt.Sprintf("%#v", ToMap(am2)) != `map[string]int{"分数":100}` {
			t.Fatal("错误")
		}
	})
}

func Test23(t *testing.T) {
	t.Run("test23 RemoveEmpty功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		RemoveEmpty(am)
		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"分数":100, "年龄":18}` {
			t.Fatal("错误")
		}
	})
}

func Test24(t *testing.T) {
	t.Run("test24 Join功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		if Join(am, ";") != `18;100;0` {
			t.Fatal("错误")
		}
	})
}

func Test25(t *testing.T) {
	t.Run("test25 JoinWithoutEmpty功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		if JoinWithoutEmpty(am, ";") != `18;100` {
			t.Fatal("错误")
		}
	})
}

func Test26(t *testing.T) {
	t.Run("test26 InByKeys功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		if !InByKeys(am, "税务") {
			t.Fatal("错误")
		}
	})
}

func Test27(t *testing.T) {
	t.Run("test27 NotInByKeys功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		if NotInByKeys(am, "税务") {
			t.Fatal("错误")
		}
	})
}

func Test28(t *testing.T) {
	t.Run("test28 InByValues功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		if !InByValues(am, 100) {
			t.Fatal("错误")
		}
	})
}

func Test29(t *testing.T) {
	t.Run("test29 NotInByValues功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		if NotInByValues(am, 18) {
			t.Fatal("错误")
		}
	})
}

func Test30(t *testing.T) {
	t.Run("test30 RemoveByKeys功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		RemoveByKeys(am, "分数", "年龄")

		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"税务":0}` {
			t.Fatal("错误")
		}
	})
}

func Test31(t *testing.T) {
	t.Run("test31 RemoveByValues功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		RemoveByValues(am, 18, 100)

		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"税务":0}` {
			t.Fatal("错误")
		}
	})
}

func Test32(t *testing.T) {
	t.Run("test32 Every功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		Every(am, func(idx int, key string, value int) (string, int) { return key, value + 1 })

		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"分数":101, "年龄":19, "税务":1}` {
			t.Fatal("错误")
		}
	})
}

func Test33(t *testing.T) {
	t.Run("test33 Each功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})
		am2 := MakeAnyMap[string, int]()

		Each(am, func(idx int, key string, value int) { Set(am2, key, value*2) })

		if fmt.Sprintf("%#v", ToMap(am2)) != `map[string]int{"分数":200, "年龄":36, "税务":0}` {
			t.Fatal("错误")
		}
	})
}

func Test34(t *testing.T) {
	t.Run("test34 Clean功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		Clean(am)

		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{}` {
			t.Fatal("错误")
		}
	})
}

func Test35(t *testing.T) {
	t.Run("test35 Cast功能", func(t *testing.T) {
		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100, "税务": 0})

		am2 := Cast[string, int, string](am, func(value int) string { return fmt.Sprintf("%v", value) })
		if fmt.Sprintf("%#v", ToMap(am2)) != `map[string]string{"分数":"100", "年龄":"18", "税务":"0"}` {
			t.Fatal("错误")
		}
	})
}
