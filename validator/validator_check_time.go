package validator

import (
	"reflect"
	"time"
)

func (my *Validator[T]) checkTime(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	if !reflect.DeepEqual(value, time.Time{}) {
		return TimeErr.NewFormat("[%s]必须是时间类型", fieldName)
	}
	return nil
}
