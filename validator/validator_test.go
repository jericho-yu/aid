package validator

import (
	"testing"
)

type TestStruct struct {
	Name  string  `v-rules:"required;min>3;max<10" v-name:"名称"`
	Email string  `v-rules:"required;email" v-name:"邮箱"`
	Date  string  `v-rules:"required;date" v-name:"日期"`
	Ptr   *string `v-rules:"required" v-name:"指针"`
}

func TestValidator(t *testing.T) {
	// 测试通过的情况
	validEmail := "test@example.com"
	validDate := "2023-10-01"
	validPtr := "valid"
	validStruct := TestStruct{
		Name:  "ValidName",
		Email: validEmail,
		Date:  validDate,
		Ptr:   &validPtr,
	}

	validator := NewValidator(validStruct)
	if err := validator.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 测试 Name 字段为空
	invalidStruct := validStruct
	invalidStruct.Name = ""
	validator = NewValidator(invalidStruct)
	if err := validator.Validate(); err == nil {
		t.Errorf("expected error for empty Name, got nil")
	}

	// 测试 Email 字段格式错误
	invalidStruct = validStruct
	invalidStruct.Email = "invalid-email"
	validator = NewValidator(invalidStruct)
	if err := validator.Validate(); err == nil {
		t.Errorf("expected error for invalid Email, got nil")
	}

	// 测试 Date 字段格式错误
	invalidStruct = validStruct
	invalidStruct.Date = "01-10-2023"
	validator = NewValidator(invalidStruct)
	if err := validator.Validate(); err == nil {
		t.Errorf("expected error for invalid Date, got nil")
	}

	// 测试 Ptr 字段为空
	invalidStruct = validStruct
	invalidStruct.Ptr = nil
	validator = NewValidator(invalidStruct)
	if err := validator.Validate(); err == nil {
		t.Errorf("expected error for nil Ptr, got nil")
	}
}
