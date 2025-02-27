package anyMap

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type AnyMap[K comparable, V any] struct {
	mu     sync.RWMutex
	keys   []K
	values []V
}

func NewAnyMap[K comparable, V any](m map[K]V) *AnyMap[K, V] {
	var (
		keys   = make([]K, 0, len(m))
		values = make([]V, 0, len(m))
	)

	for key, value := range m {
		keys = append(keys, key)
		values = append(values, value)
	}

	return &AnyMap[K, V]{
		keys:   keys,
		values: values,
		mu:     sync.RWMutex{},
	}
}

func MakeAnyMap[K comparable, V any]() *AnyMap[K, V] {
	return &AnyMap[K, V]{keys: make([]K, 0), values: make([]V, 0), mu: sync.RWMutex{}}
}

func GetKeyByIndex[K comparable, V any](am *AnyMap[K, V], idx int) K {
	var zero K
	if am == nil {
		return zero
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.keys[idx]
}

func GetKeysByIndexes[K comparable, V any](am *AnyMap[K, V], indexes ...int) []K {
	var ret = make([]K, 0)
	if am == nil {
		return ret
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, idx := range indexes {
		ret = append(ret, am.keys[idx])
	}

	return ret
}

func GetKeyByValue[K comparable, V any](am *AnyMap[K, V], v V) K {
	var zero K
	if am == nil {
		return zero
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	for idx, value := range am.values {
		if reflect.DeepEqual(v, value) {
			return am.keys[idx]
		}
	}

	return zero
}

func GetKeysByValues[K comparable, V any](am *AnyMap[K, V], values ...V) []K {
	var ret = make([]K, 0)
	if am == nil {
		return ret
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	ret = getKeysByValues(am, values...)
	return ret
}

func getKeysByValues[K comparable, V any](am *AnyMap[K, V], values ...V) []K {
	var ret = make([]K, 0)
	if am == nil {
		return ret
	}

	for idx, value := range am.values {
		for _, val := range values {
			if reflect.DeepEqual(value, val) {
				ret = append(ret, am.keys[idx])
			}
		}
	}

	return ret
}

func GetValueByIndex[K comparable, V any](am *AnyMap[K, V], idx int) V {
	var zero V
	if am == nil {
		return zero
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.values[idx]
}

func GetValuesByIndex[K comparable, V any](am *AnyMap[K, V], indexes ...int) []V {
	var ret = make([]V, 0)
	if am == nil {
		return ret
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, idx := range indexes {
		ret = append(ret, am.values[idx])
	}

	return ret
}

func GetValueByKey[K comparable, V any](am *AnyMap[K, V], k K) V {
	var zero V
	if am == nil {
		return zero
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	for idx, key := range am.keys {
		if k == key {
			return am.values[idx]
		}
	}

	return zero
}

func GetValuesByKeys[K comparable, V any](am *AnyMap[K, V], keys ...K) []V {
	var ret = make([]V, 0)
	if am == nil {
		return ret
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	for idx, key := range am.keys {
		for _, k := range keys {
			if key == k {
				ret = append(ret, am.values[idx])
			}
		}
	}

	return ret
}

func GetIndexByKey[K comparable, V any](am *AnyMap[K, V], k K) int {
	if am == nil {
		return -1
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	return getIndexByKey(am, k)
}

func getIndexByKey[K comparable, V any](am *AnyMap[K, V], k K) int {
	if am == nil {
		return -1
	}

	for idx, key := range am.keys {
		if key == k {
			return idx
		}
	}

	return -1
}

func GetIndexesByKeys[K comparable, V any](am *AnyMap[K, V], keys ...K) []int {
	indices := make([]int, 0, len(keys))
	for _, key := range keys {
		idx := getIndexByKey(am, key)
		indices = append(indices, idx)
	}
	return indices
}

func GetIndexByValue[K comparable, V any](am *AnyMap[K, V], v V) int {
	if am == nil {
		return -1
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	for idx, value := range am.values {
		if reflect.DeepEqual(value, v) {
			return idx
		}
	}

	return -1
}

func GetIndexesByValues[K comparable, V any](am *AnyMap[K, V], values ...V) []int {
	indices := make([]int, 0, len(values))
	for _, value := range values {
		idx := GetIndexByValue(am, value)
		indices = append(indices, idx)
	}
	return indices
}

func IsEmpty[K comparable, V any](am *AnyMap[K, V]) bool {
	if am == nil {
		return true
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return isEmpty(am)
}

func isEmpty[K comparable, V any](am *AnyMap[K, V]) bool {
	if am == nil {
		return true
	}

	return len(am.keys) == 0 || am.keys == nil
}

func IsNotEmpty[K comparable, V any](am *AnyMap[K, V]) bool {
	if am == nil {
		return false
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return !isEmpty(am)
}

func Set[K comparable, V any](am *AnyMap[K, V], key K, value V) {
	if am == nil {
		am = NewAnyMap(map[K]V{key: value})
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	am.keys = append(am.keys, key)
	am.values = append(am.values, value)
}

func Copy[K comparable, V any](am *AnyMap[K, V]) *AnyMap[K, V] {
	if am == nil {
		return MakeAnyMap[K, V]()
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	var m = make(map[K]V)

	for idx, key := range am.keys {
		m[key] = am.values[idx]
	}

	return NewAnyMap(m)
}

func Len[K comparable, V any](am *AnyMap[K, V]) int {
	if am == nil {
		return 0
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return len(am.keys)
}

func ToMap[K comparable, V any](am *AnyMap[K, V]) map[K]V {
	var data = make(map[K]V)

	if am == nil {
		return data
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return toMap(am)
}

func toMap[K comparable, V any](am *AnyMap[K, V]) map[K]V {
	var data = make(map[K]V)

	if am == nil {
		return data
	}

	for idx, key := range am.keys {
		data[key] = am.values[idx]
	}

	return data
}

func First[K comparable, V any](am *AnyMap[K, V]) (K, V) {
	var (
		zKey K
		zVal V
	)

	if am == nil {
		return zKey, zVal
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.keys[0], am.values[0]
}

func Last[K comparable, V any](am *AnyMap[K, V]) (K, V) {
	var (
		zKey K
		zVal V
	)

	if am == nil {
		return zKey, zVal
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.keys[len(am.keys)-1], am.values[len(am.values)-1]
}

func Filter[K comparable, V any](am *AnyMap[K, V], fn func(key K, value V) bool) *AnyMap[K, V] {
	if am == nil {
		return nil
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	return filter(am, fn)
}

func filter[K comparable, V any](am *AnyMap[K, V], fn func(key K, value V) bool) *AnyMap[K, V] {
	if am == nil {
		return nil
	}

	var data = make(map[K]V)

	for idx, key := range am.keys {
		if fn(key, am.values[idx]) {
			data[key] = am.values[idx]
		}
	}

	return NewAnyMap(data)
}

func RemoveEmpty[K comparable, V any](am *AnyMap[K, V]) {
	if am == nil {
		return
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	removeEmpty(am)
}

func removeEmpty[K comparable, V any](am *AnyMap[K, V]) {
	if am == nil {
		return
	}

	var data = filter(am, func(key K, value V) bool {
		ref := reflect.ValueOf(value)

		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				return false
			}
			if ref.Elem().IsZero() {
				return false
			}
		} else {
			if ref.IsZero() {
				return false
			}
		}
		return true
	})

	am.keys = data.keys
	am.values = data.values
}

func Join[K comparable, V any](am *AnyMap[K, V], sep string) string {
	if am == nil {
		return ""
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	return join(am, sep)
}

func join[K comparable, V any](am *AnyMap[K, V], sep string) string {
	if am == nil {
		return ""
	}

	values := make([]string, len(am.values))
	for idx, value := range am.values {
		values[idx] = fmt.Sprintf("%v", value)
	}

	return strings.Join(values, sep)
}

func JoinWithoutEmpty[K comparable, V any](am *AnyMap[K, V], sep string) string {
	if am == nil {
		return ""
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	removeEmpty(am)
	return join(am, sep)
}

func InByKeys[K comparable, V any](am *AnyMap[K, V], keys ...K) bool { return inByKeys(am, keys...) }

func inByKeys[K comparable, V any](am *AnyMap[K, V], keys ...K) bool {
	return len(GetIndexesByKeys(am, keys...)) > 0
}

func NotInByKeys[K comparable, V any](am *AnyMap[K, V], keys ...K) bool {
	return !inByKeys(am, keys...)
}

func InByValues[K comparable, V any](am *AnyMap[K, V], values ...V) bool {
	return inByValues(am, values...)
}

func inByValues[K comparable, V any](am *AnyMap[K, V], values ...V) bool {
	return len(GetIndexesByValues(am, values...)) > 0
}

func NotInByValues[K comparable, V any](am *AnyMap[K, V], values ...V) bool {
	return !inByValues(am, values...)
}

func RemoveByKeys[K comparable, V any](am *AnyMap[K, V], keys ...K) {
	if am == nil {
		return
	}

	var (
		data   = toMap(am)
		newMap *AnyMap[K, V]
	)

	am.mu.Lock()
	defer am.mu.Unlock()

	for _, key := range keys {
		delete(data, key)
	}
	newMap = NewAnyMap(data)

	am.keys = newMap.keys
	am.values = newMap.values
}

func RemoveByValues[K comparable, V any](am *AnyMap[K, V], values ...V) {
	if am == nil {
		return
	}

	var (
		data   = ToMap(am)
		newMap *AnyMap[K, V]
		keys   = getKeysByValues(am, values...)
	)

	am.mu.Lock()
	defer am.mu.Unlock()

	for _, key := range keys {
		delete(data, key)
	}
	newMap = NewAnyMap(data)

	am.keys = newMap.keys
	am.values = newMap.values
}

func Every[K comparable, V any](am *AnyMap[K, V], fn func(idx int, key K, value V) (K, V)) {
	if am == nil {
		return
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	for idx := range am.keys {
		am.keys[idx], am.values[idx] = fn(idx, am.keys[idx], am.values[idx])
	}
}

func Each[K comparable, V any](am *AnyMap[K, V], fn func(idx int, key K, value V)) {
	if am == nil {
		return
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	for idx := range am.keys {
		fn(idx, am.keys[idx], am.values[idx])
	}
}

func Clean[K comparable, V any](am *AnyMap[K, V]) {
	if am == nil {
		return
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	am.keys = make([]K, 0)
	am.values = make([]V, 0)
}

func Cast[K comparable, SRC any, DST any](am *AnyMap[K, SRC], fn func(value SRC) DST) *AnyMap[K, DST] {
	if am == nil {
		return nil
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	var data = make(map[K]DST)
	for idx, key := range am.keys {
		data[key] = fn(am.values[idx])
	}

	return NewAnyMap(data)
}
