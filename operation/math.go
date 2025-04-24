package operation

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// Min 获取最小值
func Min[T constraints.Ordered](values ...T) (T, error) {
	var zero T

	if len(values) == 0 {
		return zero, errors.New("至少需要一个值")
	}

	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min, nil
}

// Max 获取最大值
func Max[T constraints.Ordered](values ...T) (T, error) {
	var zero T

	if len(values) == 0 {
		return zero, errors.New("至少需要一个值")
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max, nil
}
