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
	return &OrderlyDict[K, V]{Key: key, Value: value}
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
func (my *AnyOrderlyDict[K, V]) SetByIndex(index int, key K, value V) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data.Set(index, NewOrderlyDict(key, value))

	return my
}

// SetByKey 设置值：根据key
func (my *AnyOrderlyDict[K, V]) SetByKey(key K, value V) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	// 如果不存在则创建
	if !my.keys.In(key) {
		my.keys.Append(key)
		my.data.Append(NewOrderlyDict(key, value))
		return my
	}

	// 如果存在则修改
	for _, val := range my.data.All() {
		if val.Key == key {
			val.Value = value
			return my
		}
	}

	return my
}

// Get 获取值
func (my *AnyOrderlyDict[K, V]) Get(key K) (V, bool) {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var emptyValue V

	for _, val := range my.data.All() {
		if val.Key == key {
			return val.Value, true
		}
	}

	return emptyValue, false
}

// Has 检查是否具备谋个key
func (my *AnyOrderlyDict[K, V]) Has(key K) bool {
	_, exist := my.Get(key)

	return exist
}

// First 获取第一个键值对
func (my *AnyOrderlyDict[K, V]) First() *OrderlyDict[K, V] { return my.data.First() }

// FirstKey 获取第一个key
func (my *AnyOrderlyDict[K, V]) FirstKey() K { return my.data.First().Key }

// FirstValue 获取第一个值
func (my *AnyOrderlyDict[K, V]) FirstValue() V { return my.data.First().Value }

// Last 获取最后一个键值对
func (my *AnyOrderlyDict[K, V]) Last() *OrderlyDict[K, V] { return my.data.Last() }

// LastKey 获取最后一个key
func (my *AnyOrderlyDict[K, V]) LastKey() K { return my.data.Last().Key }

// LastValue 获取最后一个值
func (my *AnyOrderlyDict[K, V]) LastValue() V { return my.data.Last().Value }

// Keys 获取所有key
func (my *AnyOrderlyDict[K, V]) Keys() []K { return my.keys.All() }

// ToMap 获取map格式数据
func (my *AnyOrderlyDict[K, V]) ToMap() map[K]V {
	var ret = make(map[K]V)
	for _, datum := range my.data.All() {
		ret[datum.Key] = datum.Value
	}

	return ret
}

// All 获取所有值
func (my *AnyOrderlyDict[K, V]) All() []*OrderlyDict[K, V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.data.All()
}

// Clean 清理
func (my *AnyOrderlyDict[K, V]) Clean() *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.data.Clean()
	my.keys.Clean()

	return my
}

// Len 长度
func (my *AnyOrderlyDict[K, V]) Len() int { return my.data.Len() }

// Filter 通过条件过滤
func (my *AnyOrderlyDict[K, V]) Filter(fn func(dict *OrderlyDict[K, V]) bool) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	j := 0
	ret := make([]*OrderlyDict[K, V], my.data.Len())
	for i := 0; i < my.data.Len(); i++ {
		if fn(my.data.Get(i)) {
			ret[j] = my.data.Get(i)
			j++
		}
	}

	my.data.Clean()
	my.data = array.NewAnyArray[*OrderlyDict[K, V]](ret)

	return my
}

// Copy 拷贝对象
func (my *AnyOrderlyDict[K, V]) Copy() *AnyOrderlyDict[K, V] {
	return NewAnyOrderlyDict(my.ToMap(), my.keys.All()...)
}
