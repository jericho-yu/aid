package anyMap

// func Test1(t *testing.T) {
// 	t.Run("test1 初始化一个空的AnyMap", func(t *testing.T) {
// 		am := MakeAnyMap[int, string]()
// 		if am == nil {
// 			t.Fatalf("错误")
// 		}
// 	})
// }
//
// func Test2(t *testing.T) {
// 	t.Run("test2 初始化AnyMap", func(t *testing.T) {
// 		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})
//
// 		if fmt.Sprintf("%#v", ToMap(am)) != `map[string]int{"分数":100, "年龄":18}` {
// 			t.Fatal("错误")
// 		}
// 	})
// }
//
// func Test3(t *testing.T) {
// 	t.Run("test3 GetKeyByIndex功能", func(t *testing.T) {
// 		am := NewAnyMap(map[string]int{"年龄": 18, "分数": 100})
//
// 		if GetKeyByIndex(am, 1) != "分数" {
// 			t.Fatal("错误")
// 		}
// 	})
// }
//
// func Test4(t *testing.T) {
//
// }
