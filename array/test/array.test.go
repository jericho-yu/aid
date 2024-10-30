package main

import (
	"fmt"

	"github.com/jericho-yu/aid/array"
)

type A struct {
	Name string
	Age  int
}

func main() {
	var (
		a = []A{
			{Name: "张三", Age: 18},
			{Name: "张三", Age: 19},
			{Name: "李四", Age: 18},
			{Name: "李四", Age: 19},
		}
		b any
	)
	b, _ = array.GroupBy(a, "Age")
	fmt.Printf("%+v\n", b)
}
