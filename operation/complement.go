package operation

// 取两个数组补集
func GetComplement[T any](setA, setB []T) []T {
	complement := []T{}
	elementMap := make(map[any]bool)

	// 将 setB 的元素存入 map
	for _, b := range setB {
		elementMap[b] = true
	}

	// 遍历 setA，检查是否不在 setB 中
	for _, a := range setA {
		if !elementMap[a] {
			complement = append(complement, a)
		}
	}

	return complement
}
