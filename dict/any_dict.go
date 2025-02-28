package dict

import (
	"encoding/json"
	"fmt"
	"github.com/jericho-yu/aid/array"
	"reflect"
	"strings"
	"sync"
)

type AnyDict[K comparable, V any] struct {
	data   map[K]V
	keys   []K
	values []V
	mu     sync.RWMutex
}

func MakeAnyDict[K comparable, V any]() *AnyDict[K, V] {
	return &AnyDict[K, V]{
		data:   make(map[K]V),
		keys:   make([]K, 0),
		values: make([]V, 0),
		mu:     sync.RWMutex{},
	}
}

func (my *AnyDict[K, V]) getKeyByIndex(idx int) K { return my.keys[idx] }

func (my *AnyDict[K, V]) GetKeyByIndex(idx int) K {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getKeyByIndex(idx)
}

func (my *AnyDict[K, V]) getKeysByIndexes(indexes ...int) *array.AnyArray[K] {
	keys := make([]K, 0, len(indexes))

	for _, idx := range indexes {
		keys = append(keys, my.getKeyByIndex(idx))
	}

	return array.NewAnyArray(keys)
}

func (my *AnyDict[K, V]) GetKeysByIndexes(indexes ...int) *array.AnyArray[K] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getKeysByIndexes(indexes...)
}

func (my *AnyDict[K, V]) getKeyByValue(value V) K {
	var zero K

	for key, val := range my.data {
		if reflect.DeepEqual(val, value) {
			return key
		}
	}

	return zero
}

func (my *AnyDict[K, V]) GetKeyByValue(value V) K {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getKeyByValue(value)
}

func (my *AnyDict[K, V]) getKeysByValues(values ...V) *array.AnyArray[K] {
	var ret = make([]K, 0)

	for _, value := range values {
		ret = append(ret, my.getKeyByValue(value))
	}

	return array.NewAnyArray(ret)
}

func (my *AnyDict[K, V]) GetKeysByValues(values ...V) *array.AnyArray[K] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getKeysByValues(values...)
}

func (my *AnyDict[K, V]) getValueByIndex(index int) V { return my.values[index] }

func (my *AnyDict[K, V]) GetValueByIndex(index int) V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getValueByIndex(index)
}

func (my *AnyDict[K, V]) getValuesByIndexes(indexes ...int) *array.AnyArray[V] {
	values := make([]V, 0, len(indexes))

	for _, idx := range indexes {
		values = append(values, my.values[idx])
	}

	return array.NewAnyArray(values)
}

func (my *AnyDict[K, V]) GetValuesByIndexes(indexes ...int) *array.AnyArray[V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getValuesByIndexes(indexes...)
}

func (my *AnyDict[K, V]) getValueByKey(key K) V { return my.data[key] }

func (my *AnyDict[K, V]) GetValueByKey(key K) V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getValueByKey(key)
}

func (my *AnyDict[K, V]) getValuesByKeys(keys ...K) *array.AnyArray[V] {
	values := make([]V, 0, len(keys))

	for _, key := range keys {
		values = append(values, my.data[key])
	}

	return array.NewAnyArray(values)
}

func (my *AnyDict[K, V]) GetValuesByKeys(keys ...K) *array.AnyArray[V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getValuesByKeys(keys...)
}

func (my *AnyDict[K, V]) getIndexByKey(key K) int {
	for i, k := range my.keys {
		if k == key {
			return i
		}
	}

	return -1
}

func (my *AnyDict[K, V]) GetIndexByKey(key K) int {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getIndexByKey(key)
}

func (my *AnyDict[K, V]) getIndexesByKeys(keys ...K) *array.AnyArray[int] {
	var ret = make([]int, 0)

	for _, key := range keys {
		for idx, k := range my.keys {
			if k == key {
				ret = append(ret, idx)
			}
		}
	}

	return array.NewAnyArray(ret)
}

func (my *AnyDict[K, V]) GetIndexesByKeys(keys ...K) *array.AnyArray[int] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getIndexesByKeys(keys...)
}

func (my *AnyDict[K, V]) getIndexByValue(value V) int {
	for idx, val := range my.values {
		if reflect.DeepEqual(val, value) {
			return idx
		}
	}

	return -1
}

func (my *AnyDict[K, V]) GetIndexByValue(value V) int {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getIndexByValue(value)
}

func (my *AnyDict[K, V]) getIndexesByValues(values ...V) *array.AnyArray[int] {
	var ret = make([]int, 0)

	for _, value := range values {
		for idx, val := range my.values {
			if reflect.DeepEqual(val, value) {
				ret = append(ret, idx)
			}
		}
	}

	return array.NewAnyArray(ret)
}

func (my *AnyDict[K, V]) GetIndexesByValues(values ...V) *array.AnyArray[int] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getIndexesByValues(values...)
}

func (my *AnyDict[K, V]) len() int { return len(my.keys) }

func (my *AnyDict[K, V]) Len() int {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.len()
}

func (my *AnyDict[K, V]) isEmpty() bool { return my.len() == 0 }

func (my *AnyDict[K, V]) IsEmpty() bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.isEmpty()
}

func (my *AnyDict[K, V]) set(key K, value V) *AnyDict[K, V] {
	my.data[key] = value
	my.keys = append(my.keys, key)
	my.values = append(my.values, value)
	return my
}

func (my *AnyDict[K, V]) Set(key K, value V) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.set(key, value)
}

func (my *AnyDict[K, V]) copy() *AnyDict[K, V] {
	var d = MakeAnyDict[K, V]()

	for idx, key := range my.keys {
		d.set(key, my.values[idx])
	}

	return d
}

func (my *AnyDict[K, V]) Copy() *AnyDict[K, V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.copy()
}

func (my *AnyDict[K, V]) toMap() map[K]V { return my.data }

func (my *AnyDict[K, V]) ToMap() map[K]V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.toMap()
}

func (my *AnyDict[K, V]) toString() string { return fmt.Sprintf("%v", my.data) }

func (my *AnyDict[K, V]) ToString() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.toString()
}

func (my *AnyDict[K, V]) getKeys() *array.AnyArray[K] { return array.NewAnyArray(my.keys) }

func (my *AnyDict[K, V]) GetKeys() *array.AnyArray[K] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getKeys()
}

func (my *AnyDict[K, V]) getValues() *array.AnyArray[V] { return array.NewAnyArray(my.values) }

func (my *AnyDict[K, V]) GetValues() *array.AnyArray[V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getValues()
}

func (my *AnyDict[K, V]) getIndexes() *array.AnyArray[int] {
	var ret = make([]int, 0, len(my.keys))

	for i := range my.keys {
		ret = append(ret, i)
	}

	return array.NewAnyArray(ret)
}

func (my *AnyDict[K, V]) GetIndexes() *array.AnyArray[int] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.getIndexes()
}

func (my *AnyDict[K, V]) firstKey() K { return my.keys[0] }

func (my *AnyDict[K, V]) FirstKey() K {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.firstKey()
}

func (my *AnyDict[K, V]) firstValue() V { return my.values[0] }

func (my *AnyDict[K, V]) FirstValue() V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.firstValue()
}

func (my *AnyDict[K, V]) lastKey() K {
	return my.keys[len(my.keys)-1]
}

func (my *AnyDict[K, V]) LastKey() K {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.lastKey()
}

func (my *AnyDict[K, V]) lastValue() V {
	return my.values[len(my.values)-1]
}

func (my *AnyDict[K, V]) LastValue() V {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.lastValue()
}

func (my *AnyDict[K, V]) filter(fn func(key K, value V) bool) *AnyDict[K, V] {
	var d = MakeAnyDict[K, V]()

	for key, value := range my.data {
		if fn(key, value) {
			d.set(key, value)
		}
	}

	return d
}

func (my *AnyDict[K, V]) Filter(fn func(key K, value V) bool) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	d := my.filter(fn)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) removeByKey(key K) *AnyDict[K, V] {
	var d = MakeAnyDict[K, V]()

	for idx, k := range my.keys {
		if k == key {
			continue
		}

		d.set(k, my.values[idx])
	}

	return d
}

func (my *AnyDict[K, V]) RemoveByKey(key K) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	d := my.removeByKey(key)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) removeByValue(value V) *AnyDict[K, V] {
	var d = MakeAnyDict[K, V]()

	for idx, val := range my.values {
		if reflect.DeepEqual(val, value) {
			continue
		}

		d.set(my.keys[idx], val)
	}

	return d
}

func (my *AnyDict[K, V]) RemoveByValue(value V) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	d := my.removeByValue(value)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) removeEmpty() *AnyDict[K, V] {
	d := MakeAnyDict[K, V]()

	for idx, value := range my.values {
		ref := reflect.ValueOf(value)

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

		d.set(my.keys[idx], value)
	}

	return d
}

func (my *AnyDict[K, V]) RemoveEmpty() *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	d := my.removeEmpty()
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) join(sep string) string {
	var values = make([]string, 0, len(my.values))

	for _, value := range my.values {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return strings.Join(values, sep)
}

func (my *AnyDict[K, V]) Join(sep string) string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.join(sep)
}

func (my *AnyDict[K, V]) joinWithoutEmpty(sep string) string {
	return my.copy().removeEmpty().join(sep)
}

func (my *AnyDict[K, V]) JoinWithoutEmpty(sep string) string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.joinWithoutEmpty(sep)
}

func (my *AnyDict[K, V]) inKeys(keys ...K) bool {
	return my.getIndexesByKeys(keys...).Len() == len(keys)
}

func (my *AnyDict[K, V]) InKeys(keys ...K) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.inKeys(keys...)
}

func (my *AnyDict[K, V]) NotInKeys(keys ...K) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return !my.inKeys(keys...)
}

func (my *AnyDict[K, V]) inValues(values ...V) bool {
	return my.getIndexesByValues(values...).Len() == len(values)
}

func (my *AnyDict[K, V]) InValues(values ...V) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.inValues(values...)
}

func (my *AnyDict[K, V]) NotInValues(values ...V) bool {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return !my.inValues(values...)
}

func (my *AnyDict[K, V]) every(fn func(key K, value V) (K, V)) *AnyDict[K, V] {
	var d = MakeAnyDict[K, V]()

	for key, value := range my.data {
		k, v := fn(key, value)
		d.set(k, v)
	}

	return d
}

func (my *AnyDict[K, V]) Every(fn func(key K, value V) (K, V)) *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	d := my.every(fn)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) each(fn func(key K, value V)) *AnyDict[K, V] {
	for key, value := range my.data {
		fn(key, value)
	}
	return my
}

func (my *AnyDict[K, V]) Each(fn func(key K, value V)) *AnyDict[K, V] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.each(fn)
}

func (my *AnyDict[K, V]) clean() *AnyDict[K, V] {
	my.data = make(map[K]V)
	my.keys = make([]K, 0)
	my.values = make([]V, 0)
	return my
}

func (my *AnyDict[K, V]) Clean() *AnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.clean()
}

func (my *AnyDict[K, V]) marshalJson() ([]byte, error) { return json.Marshal(&my.data) }

func (my *AnyDict[K, V]) MarshalJSON() ([]byte, error) {
	my.mu.Lock()
	defer my.mu.Unlock()

	return my.marshalJson()
}

func (my *AnyDict[K, V]) unmarshalJson(data []byte) error { return json.Unmarshal(data, &my.data) }

func (my *AnyDict[K, V]) UnmarshalJSON(data []byte) error {
	my.mu.Lock()
	defer my.mu.Unlock()

	if err := my.unmarshalJson(data); err != nil {
		return err
	}

	my.keys = make([]K, 0)
	my.values = make([]V, 0)
	for key, value := range my.data {
		my.keys = append(my.keys, key)
		my.values = append(my.values, value)
	}

	return nil
}

func Cast[K comparable, SRC, DST any](src *AnyDict[K, SRC], fn func(key K, value SRC) DST) *AnyDict[K, DST] {
	var d = MakeAnyDict[K, DST]()

	for key, value := range src.data {
		d.set(key, fn(key, value))
	}

	return d
}
