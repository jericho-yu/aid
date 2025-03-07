package operation

// Ternary 三元运算：通过值
func Ternary[V any](condition bool, T, F V) V {
	if condition {
		return T
	} else {
		return F
	}
}

// TernaryFunc 三元运算：通过回调函数
func TernaryFunc[V any](condition func() bool, T V, F V) V { return Ternary(condition(), T, F) }

// TernaryFuncAll 三元运算：通过回调函数，返回值也使用回调函数
func TernaryFuncAll[V any](condition func() bool, trueFn func() V, falseFn func() V) V {
	return Ternary(condition(), trueFn(), falseFn())
}
