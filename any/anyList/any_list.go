package anyList

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"math/rand"
)

type AnyList[T any] struct {
	mu   sync.RWMutex
	data []T
}

func NewAnyList[T any](list []T) *AnyList[T] { return &AnyList[T]{data: list, mu: sync.RWMutex{}} }

func MakeAnyList[T any](len int) *AnyList[T] {
	return &AnyList[T]{data: make([]T, len), mu: sync.RWMutex{}}
}

func Has[T any](list *AnyList[T], k int) bool {
	if list == nil {
		return false
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	return k >= 0 && k < len(list.data)
}

func has[T any](list *AnyList[T], k int) bool {
	if list == nil {
		return false
	}

	return k >= 0 && k < len(list.data)
}

func Set[T any](list *AnyList[T], k int, v T) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	list.data[k] = v
}

func Get[T any](list *AnyList[T], k int) (T, bool) {
	var zero T

	if list == nil {
		return zero, false
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	if !has(list, k) {
		return zero, false
	}

	return list.data[k], true
}

func ToSlice[T any](list *AnyList[T]) []T {
	if list == nil {
		return []T{}
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	return list.data
}

func GetByIndexes[T any](list *AnyList[T], keys ...int) []T {
	if list == nil {
		return []T{}
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	var ret = make([]T, 0)

	for _, key := range keys {
		if has(list, key) {
			ret = append(ret, list.data[key])
		}
	}

	return ret
}

func Append[T any](list *AnyList[T], v ...T) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	list.data = append(list.data, v...)
}

func First[T any](list *AnyList[T]) T {
	var zero T

	if list == nil {
		return zero
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	if len(list.data) == 0 {
		return zero
	}
	return list.data[0]
}

func Last[T any](list *AnyList[T]) T {
	var zero T

	if list == nil {
		return zero
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	if len(list.data) == 0 {
		return zero
	}

	return list.data[len(list.data)-1]
}

func FindIndex[T any](list *AnyList[T], comparator func(val T) bool) int {
	if list == nil {
		return -1
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	for idx, v := range list.data {
		if comparator(v) {
			return idx
		}
	}

	return -1
}

func FindIndexes[T any](list *AnyList[T], comparator func(val T) bool) []int {
	var ret []int

	if list == nil {
		return ret
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	for idx, v := range list.data {
		if comparator(v) {
			ret = append(ret, idx)
		}
	}

	return ret
}

func Copy[T any](list *AnyList[T]) *AnyList[T] {
	if list == nil {
		return MakeAnyList[T](0)
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	return NewAnyList(list.data)
}

func Shuffle[T any](list *AnyList[T]) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	randStr := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range list.data {
		j := randStr.Intn(i + 1)                                // 生成 [0, i] 范围内的随机数
		list.data[i], list.data[j] = list.data[j], list.data[i] // 交换元素
	}
}

func Len[T any](list *AnyList[T]) int {
	if list == nil {
		return 0
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	return len(list.data)
}

func Filter[T any](list *AnyList[T], fn func(v T) bool) *AnyList[T] {
	if list == nil {
		return nil
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	j := 0
	ret := make([]T, len(list.data))
	for i := range list.data {
		if fn(list.data[i]) {
			ret[j] = list.data[i]
			j++
		}
	}

	return NewAnyList(ret[:j])
}

func RemoveEmpty[T any](list *AnyList[T]) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	var data []T = make([]T, 0)

	for _, item := range list.data {
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

	list.data = data
}

func Join[T any](list *AnyList[T], sep string) string {
	if list == nil {
		return ""
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	values := make([]string, len(list.data))
	for idx, datum := range list.data {
		values[idx] = fmt.Sprintf("%v", datum)
	}

	return strings.Join(values, sep)
}

func JoinWithoutEmpty[T any](list *AnyList[T], sep string) string {
	list2 := Copy(list)
	RemoveEmpty(list2)
	return Join(list2, sep)
}

func In[T any](list *AnyList[T], target T) bool {
	if list == nil {
		return false
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	for _, element := range list.data {
		if reflect.DeepEqual(target, element) {
			return true
		}
	}

	return false
}

func NotIn[T any](list *AnyList[T], target T) bool { return !In(list, target) }

func AllEmpty[T any](list *AnyList[T]) bool {
	if list == nil {
		return false
	}

	list2 := Copy(list)
	RemoveEmpty(list2)
	return len(list2.data) == 0
}

func AllNotEmpty[T any](list *AnyList[T]) bool { return !AllEmpty(list) }

func Chunk[T any](list *AnyList[T], chunkSize int) [][]T {
	if list == nil {
		return [][]T{}
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	var chunks [][]T
	for i := 0; i < len(list.data); i += chunkSize {
		end := i + chunkSize
		if end > len(list.data) {
			end = len(list.data)
		}
		chunks = append(chunks, list.data[i:end])
	}

	return chunks
}

func Pluck[SRC, DST any](list *AnyList[SRC], fn func(item SRC) DST) []DST {
	if list == nil {
		return []DST{}
	}

	list.mu.RLock()
	defer list.mu.RUnlock()

	var ret = make([]DST, 0)
	for _, v := range list.data {
		ret = append(ret, fn(v))
	}

	return ret
}

func Unique[T any](list *AnyList[T]) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	seen := make(map[string]struct{})
	result := make([]T, 0)

	for _, value := range list.data {
		key := fmt.Sprintf("%v", value)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, value)
		}
	}

	list.data = result
}

func RemoveByIndexes[T any](list *AnyList[T], indexes ...int) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	keysMap := make(map[int]T)
	newData := make([]T, len(list.data)-len(indexes))
	i := 0

	for idx, datum := range list.data {
		keysMap[idx] = datum
	}

	for _, idx := range indexes {
		delete(keysMap, idx)
	}

	for idx := range list.data {
		if datum, ok := keysMap[idx]; ok {
			newData[i] = datum
			i++
		}
	}

	list.data = newData
}

func Every[T any](list *AnyList[T], fn func(idx int, item T) T) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	for idx := range list.data {
		list.data[idx] = fn(idx, list.data[idx])
	}
}

func Each[T any](list *AnyList[T], fn func(idx int, item T)) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	for idx := range list.data {
		fn(idx, list.data[idx])
	}
}

func Clean[T any](list *AnyList[T]) {
	if list == nil {
		return
	}

	list.mu.Lock()
	defer list.mu.Unlock()

	list.data = make([]T, 0)
}

func Cast[SRC any, DST any](list *AnyList[SRC], fn func(value SRC) DST) *AnyList[DST] {
	if list == nil {
		return MakeAnyList[DST](0)
	}

	var ret = make([]DST, len(list.data))

	list.mu.RLock()
	defer list.mu.RUnlock()

	for idx, datum := range list.data {
		ret[idx] = fn(datum)
	}

	return NewAnyList(ret)
}
