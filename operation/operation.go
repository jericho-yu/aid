package operation

// Ternary 三元运算：通过值
func Ternary[V any](condition bool, T any, F any) V {
	if condition {
		return T.(V)
	} else {
		return F.(V)
	}
}

// TernaryFunc 三元运算：通过回调函数
func TernaryFunc[V any](condition func() bool, T any, F any) V { return Ternary[V](condition(), T, F) }

// TernaryFuncAll 三元运算：通过回调函数，返回值也使用回调函数
func TernaryFuncAll[V any](condition func() bool, trueFn func() V, falseFn func() V) V {
	return Ternary[V](condition(), trueFn(), falseFn())
}
