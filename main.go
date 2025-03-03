package main

import (
	"errors"
	"fmt"
	"reflect"
)

type (
	MyError struct{ msg string }
	// MyError1 定义一个自定义错误类型
	MyError1 struct{ MyError }
	MyError2 struct{ MyError }
)

var (
	MyErr1 MyError1
	MyErr2 MyError2
)

func (my *MyError1) New(msg string) *MyError1 { return &MyError1{MyError{msg: msg}} }

func (my *MyError2) New(msg string) *MyError2 { return &MyError2{MyError{msg: msg}} }

// func (my *MyError) Error() string { return my.msg }

func (my *MyError1) Error() string { return my.msg }

func (my *MyError2) Error() string { return my.msg }

// Is 实现 Is 方法
func (my *MyError1) Is(target error) bool { return reflect.DeepEqual(target, &MyError1{}) }

func (my *MyError2) Is(target error) bool { return reflect.DeepEqual(target, &MyError2{}) }

func main() {
	err1 := MyErr1.New("Some error occurred")
	err2 := MyErr2.New("Some error occurred2")

	// 使用 errors.Is 来判断错误是否是 ErrMyError
	if errors.Is(err1, &MyError1{}) {
		fmt.Printf("Is OK1: %s\n", err1)
	} else {
		fmt.Println("Is NO1: Error not matched")
	}

	if errors.Is(err2, &MyError2{}) {
		fmt.Printf("Is OK2：%s\n", err2)
	} else {
		fmt.Println("Is NO2")
	}

	var (
		as1 *MyError1
		as2 *MyError2
	)

	if errors.As(err1, &as1) {
		fmt.Println("OK2")
	} else {
		fmt.Println("NO2")
	}

	if errors.As(err2, &as2) {
		fmt.Println("OK2")
	} else {
		fmt.Println("NO2")
	}
}
