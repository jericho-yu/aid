package validator

import (
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jericho-yu/aid/common"
)

// checkString 验证：string -> 支持的规则 required、email、email=、date、date=、time、time=、datetime、datetime=、size<、size<=、size>、size>=、range=、length=
func (my *ValidatorApp[T]) checkString(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if rule == "required" && reflect.ValueOf(value).IsNil() {
			return RequiredErr.New(fieldName)
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
		min := strings.TrimPrefix(rule, "size<=")
		if !(utf8.RuneCountInString(value.(string)) <= common.ToInt(min)) {
			return LengthErr.NewFormat("[%s]长度必须小于等于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size<"):
		min := strings.TrimPrefix(rule, "size<")
		if !(utf8.RuneCountInString(value.(string)) < common.ToInt(min)) {
			return LengthErr.NewFormat("[%s]长度必须小于[%d]", fieldName, common.ToInt(min))
		}
	case strings.HasPrefix(rule, "size>="):
		max := strings.TrimPrefix(rule, "size>=")
		if !(utf8.RuneCountInString(value.(string)) >= common.ToInt(max)) {
			return LengthErr.NewFormat("[%s]长度必须大于等于[%d]", fieldName, common.ToInt(max))
		}
	case strings.HasPrefix(rule, "size>"):
		max := strings.TrimPrefix(rule, "size>")
		if !(utf8.RuneCountInString(value.(string)) > common.ToInt(max)) {
			return LengthErr.NewFormat("[%s]长度必须大于[%d]", fieldName, common.ToInt(max))
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
		betweens := strings.Split(between, ",")
		if len(betweens) != 2 {
			return RuleErr.NewFormat("[%s]规则定义错误，规则定义错误，规则格式[d,d]", fieldName)
		}
		min := common.ToInt(betweens[0])
		max := common.ToInt(betweens[1])
		if utf8.RuneCountInString(value.(string)) < min || utf8.RuneCountInString(value.(string)) > max {
			return LengthErr.NewFormat("[%s]长度必须在[%d~%d]之间", fieldName, min, max)
		}
	}

	return nil
}
