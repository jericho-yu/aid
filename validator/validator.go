package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jericho-yu/aid/operation"
	"github.com/jericho-yu/aid/str"
)

type (
	// ValidatorApp 验证器 验证规则 -> [required] [email|datetime|date|time] [size<|size<=|size>|size>=|size=|size!=] [range=]
	ValidatorApp[T any] struct {
		data           T
		prefixNames    []string
		err            error
		emailFormat    string
		dateFormat     string
		timeFormat     string
		datetimeFormat string
		checkFunctions checkFunMap
	}

	ValidatorExCheckerApp struct{ ExFunMap exFunMap }

	checkFun    func(rule string, fieldName string, value any) error
	checkFunMap map[string]checkFun
	exFun       func(value any) error
	exFunMap    map[string]exFun
)

var (
	validatorExCheckerOnce sync.Once
	validatorExCheckerIns  *ValidatorExCheckerApp
	ValidatorExChecker     ValidatorExCheckerApp
)

// Once 单利化：额外验证器
func (*ValidatorExCheckerApp) Once() *ValidatorExCheckerApp {
	validatorExCheckerOnce.Do(func() {
		validatorExCheckerIns = &ValidatorExCheckerApp{
			ExFunMap: make(exFunMap),
		}
	})

	return validatorExCheckerIns
}

// RegisterExFun 注册额外验证函数
func (my *ValidatorExCheckerApp) RegisterExFun(name string, exFun exFun) *ValidatorExCheckerApp {
	my.ExFunMap[name] = exFun

	return my
}

// New 实例化：验证器
func New[T any](data T, prefixNames ...string) *ValidatorApp[T] {
	p := make([]string, 0)
	if len(prefixNames) > 0 {
		p = prefixNames
	}

	ins := &ValidatorApp[T]{
		data:           data,
		prefixNames:    p,
		emailFormat:    `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		dateFormat:     `^\d{4}-\d{2}-\d{2}$`,
		timeFormat:     `^\d{2}:\d{2}:\d{2}\.{0,1}\d+$`,
		datetimeFormat: `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`,
		checkFunctions: make(checkFunMap, 0),
	}

	ins.checkFunctions = checkFunMap{
		"string":     ins.checkString,
		"*string":    ins.checkString,
		"int":        ins.checkInt,
		"*int":       ins.checkInt,
		"int8":       ins.checkInt8,
		"*int8":      ins.checkInt8,
		"int16":      ins.checkInt16,
		"*int16":     ins.checkInt16,
		"int32":      ins.checkInt32,
		"*int32":     ins.checkInt32,
		"int64":      ins.checkInt64,
		"*int64":     ins.checkInt64,
		"uint":       ins.checkUint,
		"*uint":      ins.checkUint,
		"uint8":      ins.checkUint8,
		"*uint8":     ins.checkUint8,
		"uint16":     ins.checkUint16,
		"*uint16":    ins.checkUint16,
		"uint32":     ins.checkUint32,
		"*uint32":    ins.checkUint32,
		"uint64":     ins.checkUint64,
		"*uint64":    ins.checkUint64,
		"float32":    ins.checkFloat32,
		"*float32":   ins.checkFloat32,
		"float64":    ins.checkFloat64,
		"*float64":   ins.checkFloat64,
		"time.Time":  ins.checkTime,
		"*time.Time": ins.checkTime,
	}

	return ins
}

// Validate 执行验证
func (my *ValidatorApp[T]) Validate(funcs ...func() error) error {
	defer my.clean()

	if my.err != nil {
		return my.err
	}

	my.err = my.validate(my.data)
	if my.err != nil {
		return my.err
	}
	if len(funcs) > 0 {
		for _, fn := range funcs {
			if err := fn(); err != nil {
				return err
			}
		}
	}

	return my.err
}

// EmailFormat 设置email默认规则
func (my *ValidatorApp[T]) EmailFormat(emailFormat string) *ValidatorApp[T] {
	my.emailFormat = emailFormat

	return my
}

// DateFormat 设置日期默认规则
func (my *ValidatorApp[T]) DateFormat(dateFormat string) *ValidatorApp[T] {
	my.dateFormat = dateFormat

	return my
}

// TimeFormat 设置时间默认规则
func (my *ValidatorApp[T]) TimeFormat(timeFormat string) *ValidatorApp[T] {
	my.timeFormat = timeFormat

	return my
}

// DatetimeFormat 设置日期+时间默认规则
func (my *ValidatorApp[T]) DatetimeFormat(datetimeFormat string) *ValidatorApp[T] {
	my.datetimeFormat = datetimeFormat

	return my
}

func (my *ValidatorApp[T]) clean() { my.err = nil }

// validate 执行验证
func (my *ValidatorApp[T]) validate(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Ptr {
		return ValidateErr.New("不符合结构或指针")
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := range val.NumField() {
		field := val.Type().Field(i)
		if field.Anonymous {
			// 递归验证嵌套字段
			if err := New(val.Field(i).Interface(), my.prefixNames...).Validate(); err != nil {
				return err
			}
			continue
		}

		tag := field.Tag.Get("v-rule")
		if tag == "" || tag == "-" {
			continue
		}

		fieldName := my.concatFieldName(operation.Ternary(field.Tag.Get("v-name") != "", field.Tag.Get("v-name"), str.NewTransfer(val.Type().Name()).PascalToCamel()))

		for _, rule := range strings.Split(tag, ";") {
			if fn, exist := my.checkFunctions[fmt.Sprintf("%v", reflect.ValueOf(val.Field(i).Interface()).Type())]; exist {
				if err := fn(rule, fieldName, val.Field(i).Interface()); err != nil {
					return err
				}
			}
		}

		exTag := field.Tag.Get("v-ex")
		if exTag == "" || exTag == "-" {
			continue
		}

		for _, exRule := range strings.Split(exTag, ";") {
			if exFun, exist := ValidatorExChecker.Once().ExFunMap[exRule]; exist {
				if err := exFun(val.Field(i).Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (my *ValidatorApp[T]) concatFieldName(fieldName string) string {
	var concatFieldNames = make([]string, len(my.prefixNames)+1)

	if len(my.prefixNames) > 0 {
		copy(concatFieldNames, my.prefixNames)
		concatFieldNames[len(my.prefixNames)] = fieldName

		return strings.Join(concatFieldNames, ".")
	}

	return fieldName
}
