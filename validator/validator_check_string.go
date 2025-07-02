package validator

import (
	"github.com/jericho-yu/aid/array"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jericho-yu/aid/common"
)

// checkString 验证：string -> 支持的规则 required、email、email=、date、date=、time、time=、datetime、datetime=、size<、size<=、size>、size>=、range=、length=
func (my *ValidatorApp[T]) checkString(rule, fieldName string, value any) error {
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

	if value.(string) == "" {
		return nil
	}

	switch {
	case rule == "required":
		if value == "" {
			return RequiredErr.New(fieldName)
		}
	case rule == "email=":
		emailFormat := strings.TrimPrefix(rule, "email=")
		if matched, _ := regexp.MatchString(emailFormat, value.(string)); !matched {
			return EmailErr.New(fieldName)
		}
	case rule == "email":
		if matched, _ := regexp.MatchString(my.emailFormat, value.(string)); !matched {
			return EmailErr.New(fieldName)
		}
	case strings.HasPrefix(rule, "time"):
		if matched, _ := regexp.MatchString(my.timeFormat, value.(string)); !matched {
			return TimeErr.NewFormat("[%s]时间格式错误，正确格式[%s]", fieldName, my.timeFormat)
		}
	case strings.HasPrefix(rule, "time="):
		timeFormat := strings.TrimPrefix(rule, "time=")
		if matched, _ := regexp.MatchString(timeFormat, value.(string)); !matched {
			return TimeErr.NewFormat("[%s]时间格式错误，正确格式[%s]", fieldName, timeFormat)
		}
	case strings.HasPrefix(rule, "datetime="):
		datetimeFormat := strings.TrimPrefix(rule, "datetime=")
		if matched, _ := regexp.MatchString(datetimeFormat, value.(string)); !matched {
			return TimeErr.NewFormat("[%s]时间格式错误，正确格式[%s]", fieldName, datetimeFormat)
		}
	case strings.HasPrefix(rule, "datetime"):
		if matched, _ := regexp.MatchString(my.datetimeFormat, value.(string)); !matched {
			return TimeErr.NewFormat("[%s]时间格式错误，正确格式[%s]", fieldName, my.datetimeFormat)
		}
	case strings.HasPrefix(rule, "date="):
		dateFormat := strings.TrimPrefix(rule, "date=")
		if matched, _ := regexp.MatchString(dateFormat, value.(string)); !matched {
			return TimeErr.NewFormat("[%s]日期格式错误，正确格式[%s]", fieldName, dateFormat)
		}
	case strings.HasPrefix(rule, "date"):
		if matched, _ := regexp.MatchString(my.dateFormat, value.(string)); !matched {
			return TimeErr.NewFormat("[%s]日期格式错误，正确格式[%s]", fieldName, my.dateFormat)
		}
	case strings.HasPrefix(rule, "size<="):
		small := strings.TrimPrefix(rule, "size<=")
		if !(utf8.RuneCountInString(value.(string)) <= common.ToInt(small)) {
			return LengthErr.NewFormat("[%s]长度必须小于等于[%d]", fieldName, common.ToInt(small))
		}
	case strings.HasPrefix(rule, "size<"):
		large := strings.TrimPrefix(rule, "size<")
		if !(utf8.RuneCountInString(value.(string)) < common.ToInt(large)) {
			return LengthErr.NewFormat("[%s]长度必须小于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>="):
		large := strings.TrimPrefix(rule, "size>=")
		if !(utf8.RuneCountInString(value.(string)) >= common.ToInt(large)) {
			return LengthErr.NewFormat("[%s]长度必须大于等于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size>"):
		large := strings.TrimPrefix(rule, "size>")
		if !(utf8.RuneCountInString(value.(string)) > common.ToInt(large)) {
			return LengthErr.NewFormat("[%s]长度必须大于[%d]", fieldName, common.ToInt(large))
		}
	case strings.HasPrefix(rule, "size="):
		size := strings.TrimPrefix(rule, "size=")
		if utf8.RuneCountInString(value.(string)) != common.ToInt(size) {
			return LengthErr.NewFormat("[%s]长度必须等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "size!="):
		size := strings.TrimPrefix(rule, "size!=")
		if utf8.RuneCountInString(value.(string)) == common.ToInt(size) {
			return LengthErr.NewFormat("[%s]长度必须不等于[%d]", fieldName, common.ToInt(size))
		}
	case strings.HasPrefix(rule, "range="):
		between := strings.TrimPrefix(rule, "range=")
		betweenRange := strings.Split(between, ",")
		if len(betweenRange) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		small := common.ToInt(betweenRange[0])
		large := common.ToInt(betweenRange[1])
		if utf8.RuneCountInString(value.(string)) < small || utf8.RuneCountInString(value.(string)) > large {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, small, large)
		}
	case strings.HasPrefix(rule, "in="):
		inValuesStr := strings.TrimPrefix(rule, "in=")
		inValuesArr := array.New(strings.Split(inValuesStr, ","))
		if !inValuesArr.In(value.(string)) {
			return ValidateErr.NewFormat("[%s]值必须在[%s]中", fieldName, inValuesStr)
		}
	case strings.HasPrefix(rule, "not in="):
		inValuesStr := strings.TrimPrefix(rule, "not in=")
		inValuesArr := array.New(strings.Split(inValuesStr, ","))
		if inValuesArr.In(value.(string)) {
			return ValidateErr.NewFormat("[%s]值不可为以下内容：[%s]", fieldName, inValuesStr)
		}
	}

	return nil
}
