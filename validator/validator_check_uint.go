package validator

import (
	"reflect"
	"strings"

	"github.com/jericho-yu/aid/common"
)

// checkUint 验证：uint -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *Validator[T]) checkUint(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(uint) <= common.ToUint(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(uint) < common.ToUint(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(uint) >= common.ToUint(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(uint) > common.ToUint(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint) != common.ToUint(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint) == common.ToUint(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d~d]", fieldName)
		}
		min := common.ToUint(betweens[0])
		max := common.ToUint(betweens[1])
		if value.(uint) < min || value.(uint) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkUint8 验证：uint8 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *Validator[T]) checkUint8(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(uint8) <= common.ToUint8(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(uint8) < common.ToUint8(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(uint8) >= common.ToUint8(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(uint8) > common.ToUint8(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint8) != common.ToUint8(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint8) == common.ToUint8(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d~d]", fieldName)
		}
		min := common.ToUint8(betweens[0])
		max := common.ToUint8(betweens[1])
		if value.(uint8) < min || value.(uint8) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkUint16 验证：uint16 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *Validator[T]) checkUint16(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(uint16) <= common.ToUint16(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(uint16) < common.ToUint16(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(uint16) >= common.ToUint16(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(uint16) > common.ToUint16(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint16) != common.ToUint16(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint16) == common.ToUint16(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d~d]", fieldName)
		}
		min := common.ToUint16(betweens[0])
		max := common.ToUint16(betweens[1])
		if value.(uint16) < min || value.(uint16) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkUint32 验证：uint32 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *Validator[T]) checkUint32(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(uint32) <= common.ToUint32(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(uint32) < common.ToUint32(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(uint32) >= common.ToUint32(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(uint32) > common.ToUint32(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint32) != common.ToUint32(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint32) == common.ToUint32(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d~d]", fieldName)
		}
		min := common.ToUint32(betweens[0])
		max := common.ToUint32(betweens[1])
		if value.(uint32) < min || value.(uint32) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}

// checkUint64 验证：uint64 -> 支持的规则 required、size<、size<=、size>、size>=、range=
func (my *Validator[T]) checkUint64(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	switch {
	case strings.HasPrefix(rule, "size<="):
		min := strings.TrimPrefix(rule, "size<=")
		if value.(uint64) <= common.ToUint64(min) {
			return LengthErr.NewFormat("[%s]长度不能小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if value.(uint64) < common.ToUint64(min) {
			return LengthErr.NewFormat("[%s]长度不能小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if value.(uint64) >= common.ToUint64(max) {
			return LengthErr.NewFormat("[%s]长度不能大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if value.(uint64) > common.ToUint64(max) {
			return LengthErr.NewFormat("[%s]长度不能大于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if value.(uint64) != common.ToUint64(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if value.(uint64) == common.ToUint64(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweens := strings.Split(between, "~")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d~d]", fieldName)
		}
		min := common.ToUint64(betweens[0])
		max := common.ToUint64(betweens[1])
		if value.(uint64) < min || value.(uint64) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}
