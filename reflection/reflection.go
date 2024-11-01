package reflection

import (
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/jericho-yu/aid/array"
	"github.com/jericho-yu/aid/operation"
	"github.com/jericho-yu/aid/str"
)

type (
	Reflection struct {
		original any
		refValue reflect.Value
		refType  reflect.Type
		IsPtr    bool
		IsZero   bool
		IsTime   bool // 是否是时间
	}

	ReflectionType string
	AnyType        string
)

const (
	Int               ReflectionType = "I"
	Int8              ReflectionType = "I8"
	Int16             ReflectionType = "I16"
	Int32             ReflectionType = "I32"
	Int64             ReflectionType = "I64"
	Uint              ReflectionType = "U"
	Uint8             ReflectionType = "U8"
	Uint16            ReflectionType = "U16"
	Uint32            ReflectionType = "U32"
	Uint64            ReflectionType = "U64"
	String            ReflectionType = "STR"
	Float32           ReflectionType = "F32"
	Float64           ReflectionType = "F64"
	Datetime          ReflectionType = "DT"
	Bool              ReflectionType = "B"
	Array             ReflectionType = "A"
	Map               ReflectionType = "M"
	Struct            ReflectionType = "S"
	Nil               ReflectionType = "N"
	PtrSliceAny       ReflectionType = "*[]ANY"
	PtrSlicePtrStruct ReflectionType = "*[]*S"
	PtrSliceStruct    ReflectionType = "*[]S"
	PtrSlicePtrMap    ReflectionType = "*[]*M"
	PtrSliceMap       ReflectionType = "*[]M"
	PtrStruct         ReflectionType = "*S"
	PtrPtrStruct      ReflectionType = "**S"
	PtrPtrMap         ReflectionType = "**M"
	PtrMap            ReflectionType = "*M" // 1
	Any               ReflectionType = "ANY"
	UnKnowType        ReflectionType = "UKT"
)

// New 实例化：反射帮助
func New(object any) *Reflection {
	var (
		ins      *Reflection = &Reflection{original: object}
		refType  reflect.Type
		refValue reflect.Value
	)

	if object == nil {
		refType = reflect.TypeOf(nil)
		refValue = reflect.ValueOf(nil)
	} else {
		refType = reflect.TypeOf(object)
		refValue = reflect.ValueOf(object)
	}

	// 如果 obj 是指针，获取其元素
	if refType == nil {
		ins.refValue = reflect.Value{}
		ins.refType = reflect.TypeOf(ins.refValue)
		ins.IsZero = true
	} else if refValue.Kind() == reflect.Ptr {
		ins.refValue = refValue.Elem()
		ins.refType = refType.Elem()
	} else {
		ins.refValue = refValue
		ins.refType = refType
	}

	ins.IsPtr = ins.refType.Kind() == reflect.Ptr // 判断是否是指针

	if ins.refType == reflect.TypeOf(reflect.Value{}) {
		ins.IsZero = true
	} else if ins.GetReflectionType() == Array || ins.GetReflectionType() == Map {
		ins.IsZero = ins.refValue.Len() == 0
	} else {
		if !ins.IsZero {
			ins.IsZero = ins.refValue.IsValid() && ins.refValue.IsZero()
			// ins.IsZero = ins.refValue.IsValid() && ins.refValue.IsZero() && !ins.refValue.IsNil()
		}
	}

	// 判断是否是时间
	if !ins.IsZero {
		if ins.refValue.IsValid() {
			ins.IsTime = ins.refValue.Type() == reflect.TypeOf(time.Time{})
		}
	}

	return ins
}

// NewByReflectValue 实例化：通过reflect.Value
func NewByReflectValue(refValue reflect.Value) *Reflection {
	return New(refValue.Interface())
}

// GetValue 获取reflect.Value
func (r *Reflection) GetValue() reflect.Value { return r.refValue }

// GetType 获取reflect.Type
func (r *Reflection) GetType() reflect.Type { return r.refType }

// GetReflectionType 获取Reflection类型
func (r *Reflection) GetReflectionType() ReflectionType {
	var ref reflect.Value = reflect.ValueOf(r.original)

	if ref.Kind() != reflect.Ptr {
		var is64 bool = unsafe.Sizeof(uintptr(0)) == 8

		if r.IsSame(time.Time{}) {
			return Datetime
		}

		switch r.GetType().Kind() {
		case reflect.Int:
			return operation.Ternary[ReflectionType](is64, Int64, Int32)
		case reflect.Int8:
			return Int8
		case reflect.Int16:
			return Int16
		case reflect.Int32:
			return Int32
		case reflect.Int64:
			return Int64
		case reflect.Uint:
			return operation.Ternary[ReflectionType](is64, Uint64, Uint32)
		case reflect.Uint8:
			return Uint8
		case reflect.Uint16:
			return Uint16
		case reflect.Uint32:
			return Uint32
		case reflect.Uint64:
			return Uint64
		case reflect.Float32:
			return Float32
		case reflect.Float64:
			return Float64
		case reflect.Bool:
			return Bool
		case reflect.String:
			return String
		case reflect.Array, reflect.Slice:
			return Array
		case reflect.Map:
			return Map
		case reflect.Struct:
			return Struct
		default:
			return Nil
		}
	}

	var (
		elem     reflect.Value
		elemType reflect.Type
	)

	elem = ref.Elem()

	if elem.Kind() != reflect.Ptr { // 如果不是指针，则判断是否是切片
		if elem.Kind() == reflect.Slice { // 如果是切片，则判断切片元素是否是指针
			elemType = elem.Type().Elem()
			switch elemType.Kind() {
			case reflect.Ptr: // *[]*struct
				if elemType.Elem().Kind() == reflect.Struct { // *[]*struct
					return PtrSlicePtrStruct
				} else if elemType.Elem().Kind() == reflect.Map { // *[]*map
					return PtrSlicePtrMap
				} else { // *[]*Other
					return UnKnowType
				}
			case reflect.Struct: // *[]struct
				return PtrSliceStruct
			case reflect.Map: // *[]map
				return PtrSliceMap
			}
		} else if elem.Kind() == reflect.Struct { // *struct
			return PtrStruct
		} else if elem.Kind() == reflect.Map { // *map
			return PtrMap
		} else {
			return Any
		}
	} else {
		elemType = elem.Type().Elem()
		if elemType.Kind() == reflect.Struct { // **struct
			return PtrPtrStruct
		} else if elemType.Kind() == reflect.Map { // **map
			return PtrPtrMap
		} else {
			return Any
		}
	}
	return UnKnowType
}

// IsSame 判断两个类型是否相等
func (r *Reflection) IsSame(value any) bool {
	return r.refType == reflect.TypeOf(value)
}

// IsSameDeepEqual 判断两个值是否相等
func (r *Reflection) IsSameDeepEqual(value any) bool {
	return reflect.DeepEqual(r.refValue.Interface(), value)
}

// CallMethodByName 通过名称调用方法
func (r *Reflection) CallMethodByName(
	methodName string,
	values ...reflect.Value,
) []reflect.Value {
	method := r.GetValue().MethodByName(methodName)
	if method.IsValid() {
		return method.Call(values)
	}

	return nil
}

// FindFieldAndFill 递归查找字段并填充
func (r *Reflection) FindFieldAndFill(
	target,
	tagTitle,
	tagField string,
	process func(val reflect.Value),
) {
	findFieldAndFill(r.original, target, tagTitle, tagField, process)
}

// findFieldAndFill 递归查找字段并填充
func findFieldAndFill(
	model any,
	target,
	tagTitle,
	tagField string,
	process func(val reflect.Value),
) {
	var (
		refValue reflect.Value = reflect.ValueOf(model)
		refType  reflect.Type  = reflect.TypeOf(model)
	)

	if refValue.Kind() == reflect.Ptr {
		refValue = refValue.Elem()
		refType = refType.Elem()
	}

	// 遍历结构体的所有字段
	for i := 0; i < refValue.NumField(); i++ {
		fieldValue := refValue.Field(i)
		fieldType := refType.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		if fieldValue.Kind() == reflect.Struct && fieldValue.Type() != reflect.TypeOf(time.Time{}) {
			// 如果是纯结构体则进入递归
			findFieldAndFill(fieldValue.Addr().Interface(), target, tagTitle, tagField, process)
		} else {
			if compareTagAndTarget(fieldType.Tag, target, tagTitle, tagField) || str.NewTransfer(fieldType.Name).PascalToSnake() == target {
				process(fieldValue)
			}
		}
	}
}

// 匹配tag和target
func compareTagAndTarget(
	tag reflect.StructTag,
	target,
	tagTitle,
	tagField string,
) bool {
	var tagValue string = tag.Get(tagTitle)

	if tagValue == "" {
		return false
	}

	if tagField != "" {
		return array.NewAnyArray[string](strings.Split(tagValue, ";")).
			Every(func(s string) string {
				t := array.NewAnyArray[string](strings.Split(s, ":"))
				return operation.Ternary[string](t.First() == tagField, t.Last(), "")
			}).
			In(target)
	} else {
		return tagValue == target
	}
}
