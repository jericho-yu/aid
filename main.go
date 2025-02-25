package main

import (
	"fmt"

	anyList "github.com/jericho-yu/aid/any/list"
)

func main() {
	al := anyList.MakeAnyList[string](5)
	{
		anyList.Set(al, 0, "abc")
		anyList.Set(al, 1, "def")
	}

	fmt.Printf("%#v\n", anyList.All(al))
}
