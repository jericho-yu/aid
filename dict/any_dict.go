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
func (r *AnyDict[K, V]) Set(key K, value V) *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
	return r
}

// Get 获取元素
func (r *AnyDict[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exist := r.data[key]
	return val, exist
}

// All 获取全部元素
func (r *AnyDict[K, V]) All() map[K]V {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data
}

// Len 获取长度
func (r *AnyDict[K, V]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}

// Filter 过滤元素
func (r *AnyDict[K, V]) Filter(fn func(V) bool) *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, val := range r.data {
		if !fn(val) {
			delete(r.data, key)
		}
	}
	return r
}

// RemoveEmpty 清除空值元素
func (r *AnyDict[K, T]) RemoveEmpty() *AnyDict[K, T] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, val := range r.data {
		ref := reflect.ValueOf(val)

		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				delete(r.data, key)
			}
			if ref.Elem().IsZero() {
				delete(r.data, key)
			}
		} else {
			if ref.IsZero() {
				delete(r.data, key)
			}
		}
	}
	return r
}

// JoinWithoutEmpty 拼接非空元素
func (r *AnyDict[K, V]) JoinWithoutEmpty(sep string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]string, r.RemoveEmpty().Len())
	j := 0
	for _, datum := range r.data {
		values[j] = fmt.Sprintf("%v", datum)
		j++
	}
	return strings.Join(values, sep)
}

// ToAnyList 转any list
func (r *AnyDict[K, V]) ToAnyList() *array.AnyArray[V] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l := array.MakeAnyArray[V](r.Len())
	j := 0
	for _, v := range r.data {
		l.Set(j, v)
		j++
	}

	return l
}

// InKey 检查key是否存在
func (r *AnyDict[K, V]) InKey(target K) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exit := r.data[target]
	return exit
}

// InVal 检查值是否存在
func (r *AnyDict[K, V]) InVal(target V) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ToAnyList().In(target)
}

// AllEmpty 检查是否全部为空
func (r *AnyDict[K, V]) AllEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ToAnyList().AllEmpty()
}

// AnyEmpty 检查是否存在空值
func (r *AnyDict[K, V]) AnyEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ToAnyList().AnyEmpty()
}

// GetKeysByValue 通过值找到所有对应的key
func (r *AnyDict[K, V]) GetKeysByValue(value *array.AnyArray[V]) *array.AnyArray[K] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l := array.MakeAnyArray[K](0)
	for key, val := range r.data {
		if value.In(val) {
			l.Append(key)
		}
	}

	return l
}

// RemoveByKey 根据key删除元素
func (r *AnyDict[K, V]) RemoveByKey(key K) *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
	return r
}

// RemoveByKeys 根据key批量删除元素
func (r *AnyDict[K, V]) RemoveByKeys(keys ...K) *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, key := range keys {
		r.RemoveByKey(key)
	}
	return r
}

// RemoveByValue 根据值删除元素
func (r *AnyDict[K, V]) RemoveByValue(value V) *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, v := range r.data {
		if reflect.DeepEqual(v, value) {
			delete(r.data, key)
		}
	}
	return r
}

// RemoveByValues 根据值批量删除元素
func (r *AnyDict[K, V]) RemoveByValues(values ...V) *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.RemoveByKeys(r.GetKeysByValue(array.NewAnyArray[V](values)).All()...)
	return r
}

// Clean 清理数据
func (r *AnyDict[K, V]) Clean() *AnyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[K]V)

	return r
}

// Keys 获取所有的key
func (r *AnyDict[K, V]) Keys() []K { return GetKeys[K, V](r.All()) }

// Values 获取所有的value
func (r *AnyDict[K, V]) Values() []V { return GetValues[K, V](r.All()) }
