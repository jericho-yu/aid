package validator

import (
	"reflect"
	"time"
)

func (my *ValidatorApp[T]) checkTime(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		isNil := reflect.ValueOf(value).IsNil()
		if rule == "required" && isNil {
			return RequiredErr.New(fieldName)
		}
		if isNil {
			return nil
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	if !reflect.DeepEqual(value, time.Time{}) {
		return TimeErr.NewFormat("[%s]必须是时间类型", fieldName)
	}
	return nil
}
