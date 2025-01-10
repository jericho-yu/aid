package validator

import (
	"testing"
)

type TestStruct struct {
	Name     string  `v-rule:"required;min>3;max<10" v-name:"名称"`
	Email    string  `v-rule:"required;email" v-name:"邮箱"`
	Date     string  `v-rule:"required;date" v-name:"日期"`
	Time     string  `v-rule:"required;date" v-name:"时间"`
	Datetime *string `v-rule:"required;datetime" v-name:"日期时间"`
	Ptr      *string `v-rule:"required" v-name:"指针"`
	EmptyPtr *string `v-rule:"" v-name:"空指针"`
}

func TestValidator(t *testing.T) {
	// 测试通过的情况
	dt := "2000-01-02 03:04:05"
	validPtr := "valid"
	validStruct := TestStruct{
		Name:     "ValidName",
		Email:    "test@example.com",
		Date:     "2022-01-02",
		Time:     "03:04:05.12345",
		Datetime: &dt,
		Ptr:      &validPtr,
	}

	validator := NewValidator(validStruct)
	if err := validator.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

}
