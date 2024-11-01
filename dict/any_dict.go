package dict

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jericho-yu/aid/array"
)

type AnyDict[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewAnyDict[K comparable, V any](dict map[K]V) *AnyDict[K, V] {
	return &AnyDict[K, V]{
		data: dict,
		mu:   sync.RWMutex{},
	}
}

func MakeAnyDict[K comparable, V any]() *AnyDict[K, V] {
	return &AnyDict[K, V]{
		data: make(map[K]V),
		mu:   sync.RWMutex{},
	}
}

// Set 设置元素
func (my *AnyDict[K, V]) Set(key K, value V) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data[key] = value
	return my
}

// Get 获取元素
func (my *AnyDict[K, V]) Get(key K) (V, bool) {
	my.mu.RLock()
	defer my.mu.RUnlock()

	val, exist := my.data[key]
	return val, exist
}

// All 获取全部元素
func (my *AnyDict[K, V]) All() map[K]V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.data
}

// Len 获取长度
func (my *AnyDict[K, V]) Len() int {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return len(my.data)
}

// Filter 过滤元素
func (my *AnyDict[K, V]) Filter(fn func(V) bool) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for key, val := range my.data {
		if !fn(val) {
			delete(my.data, key)
		}
	}
	return my
}

// RemoveEmpty 清除空值元素
func (my *AnyDict[K, T]) RemoveEmpty() *AnyDict[K, T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for key, val := range my.data {
		ref := reflect.ValueOf(val)

		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				delete(my.data, key)
			}
			if ref.Elem().IsZero() {
				delete(my.data, key)
			}
		} else {
			if ref.IsZero() {
				delete(my.data, key)
			}
		}
	}
	return my
}

// JoinWithoutEmpty 拼接非空元素
func (my *AnyDict[K, V]) JoinWithoutEmpty(sep string) string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	values := make([]string, my.RemoveEmpty().Len())
	j := 0
	for _, datum := range my.data {
		values[j] = fmt.Sprintf("%v", datum)
		j++
	}
	return strings.Join(values, sep)
}

// ToAnyList 转any list
func (my *AnyDict[K, V]) ToAnyList() *array.AnyArray[V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	l := array.MakeAnyArray[V](my.Len())
	j := 0
	for _, v := range my.data {
		l.Set(j, v)
		j++
	}

	return l
}

// InKey 检查key是否存在
func (my *AnyDict[K, V]) InKey(target K) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	_, exit := my.data[target]
	return exit
}

// InVal 检查值是否存在
func (my *AnyDict[K, V]) InVal(target V) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.ToAnyList().In(target)
}

// AllEmpty 检查是否全部为空
func (my *AnyDict[K, V]) AllEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.ToAnyList().AllEmpty()
}

// AnyEmpty 检查是否存在空值
func (my *AnyDict[K, V]) AnyEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.ToAnyList().AnyEmpty()
}

// GetKeysByValue 通过值找到所有对应的key
func (my *AnyDict[K, V]) GetKeysByValue(value *array.AnyArray[V]) *array.AnyArray[K] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	l := array.MakeAnyArray[K](0)
	for key, val := range my.data {
		if value.In(val) {
			l.Append(key)
		}
	}

	return l
}

// RemoveByKey 根据key删除元素
func (my *AnyDict[K, V]) RemoveByKey(key K) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	delete(my.data, key)
	return my
}

// RemoveByKeys 根据key批量删除元素
func (my *AnyDict[K, V]) RemoveByKeys(keys ...K) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for _, key := range keys {
		my.RemoveByKey(key)
	}
	return my
}

// RemoveByValue 根据值删除元素
func (my *AnyDict[K, V]) RemoveByValue(value V) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for key, v := range my.data {
		if reflect.DeepEqual(v, value) {
			delete(my.data, key)
		}
	}
	return my
}

// RemoveByValues 根据值批量删除元素
func (my *AnyDict[K, V]) RemoveByValues(values ...V) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.RemoveByKeys(my.GetKeysByValue(array.NewAnyArray[V](values)).All()...)
	return my
}

// Clean 清理数据
func (my *AnyDict[K, V]) Clean() *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data = make(map[K]V)

	return my
}

// Keys 获取所有的key
func (my *AnyDict[K, V]) Keys() []K { return GetKeys[K, V](my.All()) }

// Values 获取所有的value
func (my *AnyDict[K, V]) Values() []V { return GetValues[K, V](my.All()) }
