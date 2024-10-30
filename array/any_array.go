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
func (r *AnyArray[T]) IsEmpty() bool {
	return len(r.data) == 0
}

// IsNotEmpty 是否不为空
func (r *AnyArray[T]) IsNotEmpty() bool {
	return len(r.data) > 0
}

// Has 检查是否存在
func (r *AnyArray[T]) Has(k int) bool {
	return k >= 0 && k < len(r.data)
}

// Set 设置值
func (r *AnyArray[T]) Set(k int, v T) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[k] = v
	return r
}

// Get 获取值
func (r *AnyArray[T]) Get(idx int) T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data[idx]
}

// Append 追加
func (r *AnyArray[T]) Append(v ...T) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = append(r.data, v...)
	return r
}

// First 获取第一个值
func (r *AnyArray[T]) First() T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data[0]
}

// Last 获取最后一个值
func (r *AnyArray[T]) Last() T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var t T

	return operation.Ternary[T](r.Len() > 1, r.data[len(r.data)-1], operation.Ternary[T](r.Len() == 0, t, r.data[0]))
}

// All 获取全部值
func (r *AnyArray[T]) All() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var ret = make([]T, len(r.data))
	copy(ret, r.data)
	return ret
}

// GetIndexByValue 根据值获取索引下标
func (r *AnyArray[T]) GetIndexByValue(value T) int {
	for idx, val := range r.data {
		if reflect.DeepEqual(val, value) {
			return idx
		}
	}
	return -1
}

// Copy 复制自己
func (r *AnyArray[T]) Copy() *AnyArray[T] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return NewAnyArray(r.data)
}

// Shuffle 函数用于打乱切片中的元素顺序
func (r *AnyArray[T]) Shuffle() *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	randStr := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range r.data {
		j := randStr.Intn(i + 1)                    // 生成 [0, i] 范围内的随机数
		r.data[i], r.data[j] = r.data[j], r.data[i] // 交换元素
	}

	return r
}

// Len 获取数组长度
func (r *AnyArray[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}

// Filter 过滤数组值
func (r *AnyArray[T]) Filter(fn func(v T) bool) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	j := 0
	ret := make([]T, len(r.data))
	for i := 0; i < len(r.data); i++ {
		if fn(r.data[i]) {
			ret[j] = r.data[i]
			j++
		}
	}

	r.data = ret[:j]
	return r
}

// RemoveEmpty 清除空值元素
func (r *AnyArray[T]) RemoveEmpty() *AnyArray[T] {
	r.mu.Lock()
	defer r.Clean()
	defer r.mu.Unlock()

	var data []T = make([]T, 0)

	for _, item := range r.data {
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
func (r *AnyArray[T]) Join(sep string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]string, r.Len())
	for idx, datum := range r.data {
		values[idx] = fmt.Sprintf("%v", datum)
	}

	return strings.Join(values, sep)
}

// JoinWithoutEmpty 拼接非空元素
func (r *AnyArray[T]) JoinWithoutEmpty(sep string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	values := make([]string, r.Copy().RemoveEmpty().Len())
	j := 0
	for _, datum := range r.Copy().RemoveEmpty().All() {
		values[j] = fmt.Sprintf("%v", datum)
		j++
	}
	return strings.Join(values, sep)
}

// In 检查值是否存在
func (r *AnyArray[T]) In(target T) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, element := range r.data {
		if reflect.DeepEqual(target, element) {
			return true
		}
	}
	return false
}

// AllEmpty 判断当前数组是否全空
func (r *AnyArray[T]) AllEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Copy().RemoveEmpty().Len() == 0
}

// AnyEmpty 判断当前数组中是否存在空值
func (r *AnyArray[T]) AnyEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Copy().RemoveEmpty().Len() != len(r.data)
}

// Chunk 分块
func (r *AnyArray[T]) Chunk(chunkSize int) [][]T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var chunks [][]T
	for i := 0; i < len(r.data); i += chunkSize {
		end := i + chunkSize
		if end > len(r.data) {
			end = len(r.data)
		}
		chunks = append(chunks, r.data[i:end])
	}

	return chunks
}

// Pluck 获取数组中指定字段的值
func (r *AnyArray[T]) Pluck(fn func(item T) any) *AnyArray[any] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var ret = make([]any, 0)
	for _, v := range r.data {
		ret = append(ret, fn(v))
	}

	return NewAnyArray(ret)
}

// Unique 切片去重
func (r *AnyArray[T]) Unique() *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	seen := make(map[string]struct{}) // 使用空结构体作为值，因为我们只关心键
	result := make([]T, 0)

	for _, value := range r.data {
		key := fmt.Sprintf("%v", value)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, value)
		}
	}

	r.data = result
	return r
}

// RemoveByIndexes 根据索引删除元素
func (r *AnyArray[T]) RemoveByIndexes(indexes ...int) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	newData := make([]T, len(r.data)-len(indexes))
	idx := 0
	for i, v := range r.data {
		if !In(i, indexes) {
			newData[idx] = v
			idx++
		}
	}

	return NewAnyArray[T](newData)
}

// RemoveByValue 删除数组中对应的目标
func (r *AnyArray[T]) RemoveByValue(target T) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	var ret = make([]T, len(r.data))
	j := 0
	for _, value := range r.data {
		if !reflect.DeepEqual(value, target) {
			ret[j] = value
			j++
		}
	}
	r.data = ret[:j]
	return r
}

// RemoveByValues 删除数组中对应的多个目标
func (r *AnyArray[T]) RemoveByValues(targets ...T) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, target := range targets {
		r.RemoveByValue(target)
	}
	return r
}

// Every 循环处理每一个
func (r *AnyArray[T]) Every(fn func(item T) T) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx := range r.data {
		r.data[idx] = fn(r.data[idx])
	}

	return r
}

// Each 遍历数组
func (r *AnyArray[T]) Each(fn func(idx int, item T)) *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx := range r.data {
		fn(idx, r.data[idx])
	}

	return r
}

// Clean 清理数据
func (r *AnyArray[T]) Clean() *AnyArray[T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make([]T, 0)

	return r
}
