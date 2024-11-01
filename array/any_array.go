package array

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/jericho-yu/aid/operation"
)

type AnyArray[T any] struct {
	data []T
	mu   sync.RWMutex
}

// NewAnyArray 实例化
func NewAnyArray[T any](list []T) *AnyArray[T] {
	return &AnyArray[T]{
		data: list,
		mu:   sync.RWMutex{},
	}
}

// MakeAnyArray 初始化
func MakeAnyArray[T any](size int) *AnyArray[T] {
	return &AnyArray[T]{
		data: make([]T, size),
		mu:   sync.RWMutex{},
	}
}

// IsEmpty 是否为空
func (my *AnyArray[T]) IsEmpty() bool {
	return len(my.data) == 0
}

// IsNotEmpty 是否不为空
func (my *AnyArray[T]) IsNotEmpty() bool {
	return len(my.data) > 0
}

// Has 检查是否存在
func (my *AnyArray[T]) Has(k int) bool {
	return k >= 0 && k < len(my.data)
}

// Set 设置值
func (my *AnyArray[T]) Set(k int, v T) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data[k] = v
	return my
}

// Get 获取值
func (my *AnyArray[T]) Get(idx int) T {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.data[idx]
}

// Append 追加
func (my *AnyArray[T]) Append(v ...T) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data = append(my.data, v...)
	return my
}

// First 获取第一个值
func (my *AnyArray[T]) First() T {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.data[0]
}

// Last 获取最后一个值
func (my *AnyArray[T]) Last() T {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var t T

	return operation.Ternary[T](my.Len() > 1, my.data[len(my.data)-1], operation.Ternary[T](my.Len() == 0, t, my.data[0]))
}

// All 获取全部值
func (my *AnyArray[T]) All() []T {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var ret = make([]T, len(my.data))
	copy(ret, my.data)
	return ret
}

// GetIndexByValue 根据值获取索引下标
func (my *AnyArray[T]) GetIndexByValue(value T) int {
	for idx, val := range my.data {
		if reflect.DeepEqual(val, value) {
			return idx
		}
	}
	return -1
}

// Copy 复制自己
func (my *AnyArray[T]) Copy() *AnyArray[T] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return NewAnyArray(my.data)
}

// Shuffle 函数用于打乱切片中的元素顺序
func (my *AnyArray[T]) Shuffle() *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	randStr := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range my.data {
		j := randStr.Intn(i + 1)                        // 生成 [0, i] 范围内的随机数
		my.data[i], my.data[j] = my.data[j], my.data[i] // 交换元素
	}

	return my
}

// Len 获取数组长度
func (my *AnyArray[T]) Len() int {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return len(my.data)
}

// Filter 过滤数组值
func (my *AnyArray[T]) Filter(fn func(v T) bool) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	j := 0
	ret := make([]T, len(my.data))
	for i := 0; i < len(my.data); i++ {
		if fn(my.data[i]) {
			ret[j] = my.data[i]
			j++
		}
	}

	my.data = ret[:j]
	return my
}

// RemoveEmpty 清除空值元素
func (my *AnyArray[T]) RemoveEmpty() *AnyArray[T] {
	my.mu.Lock()
	defer my.Clean()
	defer my.mu.Unlock()

	var data []T = make([]T, 0)

	for _, item := range my.data {
		ref := reflect.ValueOf(item)

		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				continue
			}
			if ref.Elem().IsZero() {
				continue
			}
		} else {
			if ref.IsZero() {
				continue
			}
		}

		data = append(data, item)
	}

	return NewAnyArray[T](data)
}

// Join 拼接字符串
func (my *AnyArray[T]) Join(sep string) string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	values := make([]string, my.Len())
	for idx, datum := range my.data {
		values[idx] = fmt.Sprintf("%v", datum)
	}

	return strings.Join(values, sep)
}

// JoinWithoutEmpty 拼接非空元素
func (my *AnyArray[T]) JoinWithoutEmpty(sep string) string {
	my.mu.Lock()
	defer my.mu.Unlock()

	values := make([]string, my.Copy().RemoveEmpty().Len())
	j := 0
	for _, datum := range my.Copy().RemoveEmpty().All() {
		values[j] = fmt.Sprintf("%v", datum)
		j++
	}
	return strings.Join(values, sep)
}

// In 检查值是否存在
func (my *AnyArray[T]) In(target T) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	for _, element := range my.data {
		if reflect.DeepEqual(target, element) {
			return true
		}
	}
	return false
}

// NotIn 检查值是否不存在
func (my *AnyArray[T]) NotIn(target T) bool {
	return !my.In(target)
}

// AllEmpty 判断当前数组是否全空
func (my *AnyArray[T]) AllEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.Copy().RemoveEmpty().Len() == 0
}

// AnyEmpty 判断当前数组中是否存在空值
func (my *AnyArray[T]) AnyEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.Copy().RemoveEmpty().Len() != len(my.data)
}

// Chunk 分块
func (my *AnyArray[T]) Chunk(chunkSize int) [][]T {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var chunks [][]T
	for i := 0; i < len(my.data); i += chunkSize {
		end := i + chunkSize
		if end > len(my.data) {
			end = len(my.data)
		}
		chunks = append(chunks, my.data[i:end])
	}

	return chunks
}

// Pluck 获取数组中指定字段的值
func (my *AnyArray[T]) Pluck(fn func(item T) any) *AnyArray[any] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var ret = make([]any, 0)
	for _, v := range my.data {
		ret = append(ret, fn(v))
	}

	return NewAnyArray(ret)
}

// Unique 切片去重
func (my *AnyArray[T]) Unique() *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	seen := make(map[string]struct{}) // 使用空结构体作为值，因为我们只关心键
	result := make([]T, 0)

	for _, value := range my.data {
		key := fmt.Sprintf("%v", value)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, value)
		}
	}

	my.data = result
	return my
}

// RemoveByIndexes 根据索引删除元素
func (my *AnyArray[T]) RemoveByIndexes(indexes ...int) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	newData := make([]T, len(my.data)-len(indexes))
	idx := 0
	for i, v := range my.data {
		if !In(i, indexes) {
			newData[idx] = v
			idx++
		}
	}

	return NewAnyArray[T](newData)
}

// RemoveByValue 删除数组中对应的目标
func (my *AnyArray[T]) RemoveByValue(target T) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	var ret = make([]T, len(my.data))
	j := 0
	for _, value := range my.data {
		if !reflect.DeepEqual(value, target) {
			ret[j] = value
			j++
		}
	}
	my.data = ret[:j]
	return my
}

// RemoveByValues 删除数组中对应的多个目标
func (my *AnyArray[T]) RemoveByValues(targets ...T) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for _, target := range targets {
		my.RemoveByValue(target)
	}
	return my
}

// Every 循环处理每一个
func (my *AnyArray[T]) Every(fn func(item T) T) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for idx := range my.data {
		my.data[idx] = fn(my.data[idx])
	}

	return my
}

// Each 遍历数组
func (my *AnyArray[T]) Each(fn func(idx int, item T)) *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for idx := range my.data {
		fn(idx, my.data[idx])
	}

	return my
}

// Clean 清理数据
func (my *AnyArray[T]) Clean() *AnyArray[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data = make([]T, 0)

	return my
}
