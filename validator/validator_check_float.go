package validator

import (
	"reflect"
	"strings"

	"github.com/jericho-yu/aid/common"
)

// checkFloat32 验证：float32 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkFloat32(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(float32) <= common.ToFloat32(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%f]", fieldName, common.ToFloat32(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(float32) < common.ToFloat32(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%f]", fieldName, common.ToFloat32(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(float32) >= common.ToFloat32(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%f]", fieldName, common.ToFloat32(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(float32) > common.ToFloat32(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%f]", fieldName, common.ToFloat32(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(float32) != common.ToFloat32(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%f]", fieldName, common.ToFloat32(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(float32) == common.ToFloat32(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%f]", fieldName, common.ToFloat32(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式：f,f", fieldName)
		}
		min := common.ToFloat32(betweens[0])
		max := common.ToFloat32(betweens[1])
		if value.(float32) < min || value.(float32) > max {
			return LengthErr.NewFormat("[%s]长度必须在：%f~%f之间", fieldName, min, max)
		}
	}

	return nil
}

// checkFloat64 验证：float64 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkFloat64(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if !(value.(float64) <= common.ToFloat64(min)) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%f]", fieldName, common.ToFloat64(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if !(value.(float64) < common.ToFloat64(min)) {
			return LengthErr.NewFormat("[%s]长度不能小于[%f]", fieldName, common.ToFloat64(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if !(value.(float64) >= common.ToFloat64(max)) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%f]", fieldName, common.ToFloat64(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if !(value.(float64) > common.ToFloat64(max)) {
			return LengthErr.NewFormat("[%s]长度不能大于[%f]", fieldName, common.ToFloat64(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(float64) != common.ToFloat64(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%f]", fieldName, common.ToFloat64(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(float64) == common.ToFloat64(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%f]", fieldName, common.ToFloat64(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式：f,f", fieldName)
		}
		min := common.ToFloat64(betweens[0])
		max := common.ToFloat64(betweens[1])
		if value.(float64) < min || value.(float64) > max {
			return LengthErr.NewFormat("[%s]长度必须在：%f~%f之间", fieldName, min, max)
		}
	}

	return nil
}
