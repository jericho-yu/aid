package dict

import (
	"sync"

	"github.com/jericho-yu/aid/array"
)

type (
	OrderlyDict[K comparable, V any] struct {
		Key   K
		Value V
	}

	AnyOrderlyDict[K comparable, V any] struct {
		data *array.AnyArray[*OrderlyDict[K, V]]
		keys *array.AnyArray[K]
		mu   sync.RWMutex
	}
)

// NewOrderlyDict 实例化：有序字典项
func NewOrderlyDict[K comparable, V any](key K, value V) *OrderlyDict[K, V] {
	return &OrderlyDict[K, V]{
		Key:   key,
		Value: value,
	}
}

// NewAnyOrderlyDict 实例化：有序字典
func NewAnyOrderlyDict[K comparable, V any](m map[K]V, keys ...K) *AnyOrderlyDict[K, V] {
	anyOrderlyDict := &AnyOrderlyDict[K, V]{
		data: array.MakeAnyArray[*OrderlyDict[K, V]](len(keys)),
		mu:   sync.RWMutex{},
		keys: array.NewAnyArray[K](keys),
	}

	count := 0
	for _, key := range keys {
		anyOrderlyDict.data.Set(count, NewOrderlyDict[K, V](key, m[key]))
		count++
	}

	return anyOrderlyDict
}

// MakeAnyOrderlyDict 格式化：有序字典
func MakeAnyOrderlyDict[K comparable, V any](size int) *AnyOrderlyDict[K, V] {
	return &AnyOrderlyDict[K, V]{
		data: array.MakeAnyArray[*OrderlyDict[K, V]](size),
		mu:   sync.RWMutex{},
		keys: array.MakeAnyArray[K](size),
	}
}

// SetByIndex 设置值：根据索引
func (r *AnyOrderlyDict[K, V]) SetByIndex(index int, key K, value V) *AnyOrderlyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data.Set(index, NewOrderlyDict(key, value))

	return r
}

// SetByKey 设置值：根据key
func (r *AnyOrderlyDict[K, V]) SetByKey(key K, value V) *AnyOrderlyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 如果不存在则创建
	if !r.keys.In(key) {
		r.keys.Append(key)
		r.data.Append(NewOrderlyDict(key, value))
		return r
	}

	// 如果存在则修改
	for _, val := range r.data.All() {
		if val.Key == key {
			val.Value = value
			return r
		}
	}

	return r
}

// Get 获取值
func (r *AnyOrderlyDict[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var emptyValue V

	for _, val := range r.data.All() {
		if val.Key == key {
			return val.Value, true
		}
	}

	return emptyValue, false
}

// First 获取第一个键值对
func (r *AnyOrderlyDict[K, V]) First() *OrderlyDict[K, V] {
	return r.data.First()
}

// FirstKey 获取第一个key
func (r *AnyOrderlyDict[K, V]) FirstKey() K {
	return r.data.First().Key
}

// FirstValue 获取第一个值
func (r *AnyOrderlyDict[K, V]) FirstValue() V {
	return r.data.First().Value
}

// Last 获取最后一个键值对
func (r *AnyOrderlyDict[K, V]) Last() *OrderlyDict[K, V] {
	return r.data.Last()
}

// LastKey 获取最后一个key
func (r *AnyOrderlyDict[K, V]) LastKey() K {
	return r.data.Last().Key
}

// LastValue 获取最后一个值
func (r *AnyOrderlyDict[K, V]) LastValue() V {
	return r.data.Last().Value
}

// Keys 获取所有key
func (r *AnyOrderlyDict[K, V]) Keys() []K {
	return r.keys.All()
}

// ToMap 获取map格式数据
func (r *AnyOrderlyDict[K, V]) ToMap() map[K]V {
	var ret = make(map[K]V)
	for _, datum := range r.data.All() {
		ret[datum.Key] = datum.Value
	}

	return ret
}

// All 获取所有值
func (r *AnyOrderlyDict[K, V]) All() []*OrderlyDict[K, V] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data.All()
}

// Clean 清理
func (r *AnyOrderlyDict[K, V]) Clean() *AnyOrderlyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data.Clean()
	r.keys.Clean()

	return r
}

// Len 长度
func (r *AnyOrderlyDict[K, V]) Len() int {
	return r.data.Len()
}

// Filter 通过条件过滤
func (r *AnyOrderlyDict[K, V]) Filter(fn func(dict *OrderlyDict[K, V]) bool) *AnyOrderlyDict[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	j := 0
	ret := make([]*OrderlyDict[K, V], r.data.Len())
	for i := 0; i < r.data.Len(); i++ {
		if fn(r.data.Get(i)) {
			ret[j] = r.data.Get(i)
			j++
		}
	}

	r.data.Clean()
	r.data = array.NewAnyArray[*OrderlyDict[K, V]](ret)
	return r
}

// Copy 拷贝对象
func (r *AnyOrderlyDict[K, V]) Copy() *AnyOrderlyDict[K, V] {
	return NewAnyOrderlyDict(r.ToMap(), r.keys.All()...)
}
