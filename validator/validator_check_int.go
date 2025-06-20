package validator

import (
	"reflect"
	"strings"

	"github.com/jericho-yu/aid/common"
)

// checkInt 验证：int -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		isNil := reflect.ValueOf(value).IsNil()
		if isNil {
			if rule == "required" {
				return RequiredErr.New(fieldName)
			}
			return nil
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(int) <= common.ToInt(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(int) < common.ToInt(min) {
			return LengthErr.NewFormat("[%s]长度不能小于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(int) >= common.ToInt(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(int) > common.ToInt(max) {
			return LengthErr.NewFormat("[%s]长度不能大于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int) != common.ToInt(size) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int) == common.ToInt(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		min := common.ToInt(betweens[0])
		max := common.ToInt(betweens[1])
		if value.(int) < min || value.(int) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkInt8 验证：int8 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt8(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		isNil := reflect.ValueOf(value).IsNil()
		if isNil {
			if rule == "required" {
				return RequiredErr.New(fieldName)
			}
			return nil
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(int8) <= common.ToInt8(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(int8) < common.ToInt8(min) {
			return LengthErr.NewFormat("[%s]长度不能小于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(int8) >= common.ToInt8(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(int8) > common.ToInt8(max) {
			return LengthErr.NewFormat("[%s]长度不能大于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int8) != common.ToInt8(size) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int8) == common.ToInt8(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		min := common.ToInt8(betweens[0])
		max := common.ToInt8(betweens[1])
		if value.(int8) < min || value.(int8) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkInt16 验证：int16 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt16(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		isNil := reflect.ValueOf(value).IsNil()
		if isNil {
			if rule == "required" {
				return RequiredErr.New(fieldName)
			}
			return nil
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(int16) <= common.ToInt16(min) {
			return LengthErr.NewFormat("[%s]长度必须小于等于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(int16) < common.ToInt16(min) {
			return LengthErr.NewFormat("[%s]长度必须小于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(int16) >= common.ToInt16(max) {
			return LengthErr.NewFormat("[%s]长度必须大于等于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if !(value.(int16) > common.ToInt16(max)) {
			return LengthErr.NewFormat("[%s]长度必须大于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if !(value.(int16) != common.ToInt16(size)) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if !(value.(int16) == common.ToInt16(size)) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		min := common.ToInt16(betweens[0])
		max := common.ToInt16(betweens[1])
		if value.(int16) < min || value.(int16) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkInt32 验证：int32 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt32(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		isNil := reflect.ValueOf(value).IsNil()
		if isNil {
			if rule == "required" {
				return RequiredErr.New(fieldName)
			}
			return nil
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(int32) <= common.ToInt32(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(int32) < common.ToInt32(min) {
			return LengthErr.NewFormat("[%s]长度不能小于：[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(int32) >= common.ToInt32(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(int32) > common.ToInt32(max) {
			return LengthErr.NewFormat("[%s]长度不能大于：[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int32) != common.ToInt32(size) {
			return LengthErr.NewFormat("[%s]长度必须等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int32) == common.ToInt32(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于：[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		min := common.ToInt32(betweens[0])
		max := common.ToInt32(betweens[1])
		if value.(int32) < min || value.(int32) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkInt64 验证：int64 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *ValidatorApp[T]) checkInt64(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		isNil := reflect.ValueOf(value).IsNil()
		if isNil {
			if rule == "required" {
				return RequiredErr.New(fieldName)
			}
			return nil
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(int64) <= common.ToInt64(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(int64) < common.ToInt64(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(int64) >= common.ToInt64(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(int64) > common.ToInt64(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(int64) != common.ToInt64(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(int64) == common.ToInt64(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		min := common.ToInt64(betweens[0])
		max := common.ToInt64(betweens[1])
		if value.(int64) < min || value.(int64) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}
