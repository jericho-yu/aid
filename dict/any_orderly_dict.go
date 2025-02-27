package dict

import (
	"fmt"
	"reflect"
	"strings"
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

	IAnyOrderlyDict[K comparable, V any] interface {
		SetByIndex(index int, key K, value V) *AnyOrderlyDict[K, V]
		SetByKey(key K, value V) *AnyOrderlyDict[K, V]
		Get(key K) (V, bool)
		Has(key K) bool
		First() *OrderlyDict[K, V]
		FirstKey() K
		FirstValue() V
		Last() *OrderlyDict[K, V]
		LastKey() K
		LastValue() V
		Keys() []K
		ToMap() map[K]V
		All() []*OrderlyDict[K, V]
		Clean() *AnyOrderlyDict[K, V]
		Len() int
		Filter(fn func(dict *OrderlyDict[K, V]) bool) *AnyOrderlyDict[K, V]
		Copy() *AnyOrderlyDict[K, V]
		RemoveEmpty() *AnyOrderlyDict[K, V]
		Join(sep string) string
		JoinWithoutEmpty(sep string) string
		ToAnyArray() *array.AnyArray[V]
		ToAnyDict() *AnyDict[K, V]
		InKey(k K) bool
		InValue(v V) bool
		NotInKey(k K) bool
		NotInValue(v V) bool
		AllEmpty() bool
		AnyEmpty() bool
		Chunk(chunkSize int) [][]V
		Pluck(fn func(item V) V) *AnyOrderlyDict[K, V]
		Unique() *AnyOrderlyDict[K, V]
		RemoveByIndexes(indexes ...int) *AnyOrderlyDict[K, V]
		RemoveByKeys(keys ...K) *AnyOrderlyDict[K, V]
		RemoveByValues(values ...V) *AnyOrderlyDict[K, V]
		Append(key K, value V) *AnyOrderlyDict[K, V]
		Every(fn func(item V) V) *AnyOrderlyDict[K, V]
		Each(fn func(idx int, key K, value V)) *AnyOrderlyDict[K, V]
	}

	seenStruct[K comparable, V any] struct {
		key   K
		value V
		total uint
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
		keys: array.NewAnyArray(keys),
	}

	count := 0
	for _, key := range keys {
		anyOrderlyDict.data.Set(count, NewOrderlyDict(key, m[key]))
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
	my.data = array.NewAnyArray(ret)

	return my
}

// Copy 拷贝对象
func (my *AnyOrderlyDict[K, V]) Copy() *AnyOrderlyDict[K, V] {
	return NewAnyOrderlyDict(my.ToMap(), my.keys.All()...)
}

// RemoveEmpty 清空空值
func (my *AnyOrderlyDict[K, V]) RemoveEmpty() *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for idx := range my.keys.All() {
		ref := reflect.ValueOf(my.data.Get(idx))

		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				my.data.RemoveByIndexes(idx)
				my.keys.RemoveByIndexes(idx)
			}
			if ref.Elem().IsZero() {
				my.data.RemoveByIndexes(idx)
				my.keys.RemoveByIndexes(idx)
			}
		} else {
			if ref.IsZero() {
				my.data.RemoveByIndexes(idx)
				my.keys.RemoveByIndexes(idx)
			}
		}
	}

	return my
}

// Join 拼接字符串
func (my *AnyOrderlyDict[K, V]) Join(sep string) string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	values := make([]string, my.Len())
	for idx, datum := range my.data.All() {
		values[idx] = fmt.Sprintf("%v", datum.Value)
	}

	return strings.Join(values, sep)
}

// JoinWithoutEmpty 拼接非空字符串
func (my *AnyOrderlyDict[K, V]) JoinWithoutEmpty(sep string) string {
	values := make([]string, my.Copy().RemoveEmpty().Len())
	j := 0
	for _, datum := range my.Copy().RemoveEmpty().All() {
		values[j] = fmt.Sprintf("%v", datum.Value)
		j++
	}

	return strings.Join(values, sep)
}

// ToAnyArray 转AnyArray
func (my *AnyOrderlyDict[K, V]) ToAnyArray() *array.AnyArray[V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var arr = make([]V, my.Len())

	for idx, datum := range my.data.All() {
		arr[idx] = datum.Value
	}

	return array.NewAnyArray(arr)
}

// ToAnyDict 转AndDict
func (my *AnyOrderlyDict[K, V]) ToAnyDict() *AnyDict[K, V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var m = make(map[K]V)

	for _, datum := range my.data.All() {
		m[datum.Key] = datum.Value
	}

	return NewAnyDict(m)
}

// InKey 检查key是否存在
func (my *AnyOrderlyDict[K, V]) InKey(k K) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	for _, datum := range my.data.All() {
		if datum.Key == k {
			return true
		}
	}

	return false
}

// InValue 检查value是否存在
func (my *AnyOrderlyDict[K, V]) InValue(v V) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	for _, datum := range my.data.All() {
		if reflect.DeepEqual(datum.Value, v) {
			return true
		}
	}

	return false
}

// NotInKey 检查key是否不存在
func (my *AnyOrderlyDict[K, V]) NotInKey(k K) bool { return !my.InKey(k) }

// NotInValue 检查value是否不存在
func (my *AnyOrderlyDict[K, V]) NotInValue(v V) bool { return !my.InValue(v) }

// AllEmpty 检查是否非0值
func (my *AnyOrderlyDict[K, V]) AllEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.Copy().RemoveEmpty().Len() == 0
}

// AnyEmpty 判断当前数组中是否存在空值
func (my *AnyOrderlyDict[K, V]) AnyEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.Copy().RemoveEmpty().Len() != len(my.data.All())
}

// Chunk 分块
func (my *AnyOrderlyDict[K, V]) Chunk(chunkSize int) [][]V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var chunks [][]V
	for i := 0; i < len(my.data.All()); i += chunkSize {
		end := i + chunkSize

		if end > my.data.Len() {
			end = my.data.Len()
		}

		chunks = append(chunks, my.ToAnyArray().All()[i:end])
	}

	return chunks
}

// Pluck 获取数组中指定字段的值
func (my *AnyOrderlyDict[K, V]) Pluck(fn func(item V) V) *AnyOrderlyDict[K, V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	var ret = make(map[K]V)
	for _, v := range my.data.All() {
		ret[v.Key] = fn(v.Value)
	}

	return NewAnyOrderlyDict(ret)
}

// Unique 切片去重
func (my *AnyOrderlyDict[K, V]) Unique() *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	seen := make(map[string]seenStruct[K, V])
	for _, datum := range my.data.All() {
		key := fmt.Sprintf("%v", datum.Value)
		if _, exists := seen[key]; !exists {
			seen[key] = seenStruct[K, V]{
				key:   datum.Key,
				value: datum.Value,
				total: 0,
			}
		}
		temp := seen[key]
		temp.total += 1
		seen[key] = temp
	}

	result := make(map[K]V, len(seen))
	for _, datum := range seen {
		result[datum.key] = datum.value
	}

	return NewAnyOrderlyDict(result)
}

// RemoveByIndexes 根据索引移除
func (my *AnyOrderlyDict[K, V]) RemoveByIndexes(indexes ...int) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.keys.RemoveByIndexes(indexes...)
	my.data.RemoveByIndexes(indexes...)

	return my
}

// RemoveByKeys 通过key删除内容
func (my *AnyOrderlyDict[K, V]) RemoveByKeys(keys ...K) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	indexes := make([]int, 0)

	for _, key := range keys {
		for idx, datum := range my.data.All() {
			if datum.Key == key {
				indexes = append(indexes, idx)
			}
		}
	}

	my.keys.RemoveByIndexes(indexes...)
	my.data.RemoveByIndexes(indexes...)

	return my
}

// RemoveByValues 通过value删除内容
func (my *AnyOrderlyDict[K, V]) RemoveByValues(values ...V) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	indexes := make([]int, 0)

	for _, value := range values {
		for idx, datum := range my.data.All() {
			if reflect.DeepEqual(value, datum.Value) {
				indexes = append(indexes, idx)
			}
		}
	}

	my.keys.RemoveByIndexes(indexes...)
	my.data.RemoveByIndexes(indexes...)

	return my
}

// Append 追加内容
func (my *AnyOrderlyDict[K, V]) Append(key K, value V) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.keys.Append(key)
	my.data.Append(NewOrderlyDict(key, value))

	return my
}

// Every 循环处理每一个
func (my *AnyOrderlyDict[K, V]) Every(fn func(item V) V) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for idx := range my.data.All() {
		my.data.All()[idx] = &OrderlyDict[K, V]{Key: my.data.All()[idx].Key, Value: fn(my.data.All()[idx].Value)}
	}

	return my
}

// Each 遍历数组
func (my *AnyOrderlyDict[K, V]) Each(fn func(idx int, key K, value V)) *AnyOrderlyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for idx := range my.data.All() {
		fn(idx, my.data.All()[idx].Key, my.data.All()[idx].Value)
	}

	return my
}

// CastOrderlyDict 转换格式
func CastOrderlyDict[K comparable, SRC, DST any](ad *AnyOrderlyDict[K, SRC], fn func(value SRC) DST) *AnyOrderlyDict[K, DST] {
	ad.mu.RLock()
	defer ad.mu.RUnlock()

	var ret = make(map[K]DST)
	for _, datum := range ad.data.All() {
		ret[datum.Key] = fn(datum.Value)
	}

	return NewAnyOrderlyDict(ret)
}
