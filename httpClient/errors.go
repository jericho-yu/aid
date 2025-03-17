package httpClient

import (
	"fmt"
	"reflect"

	"github.com/jericho-yu/aid/array"
	"github.com/jericho-yu/aid/myError"
	"github.com/jericho-yu/aid/operation"
)

type (
	ReadResponseError    struct{ myError.MyError }
	UrlEmptyError        struct{ myError.MyError }
	GenerateCertError    struct{ myError.MyError }
	GenerateRequestError struct{ myError.MyError }
	UnmarshalXmlError    struct{ myError.MyError }
	UnmarshalJsonError   struct{ myError.MyError }
	SetSteamBodyError    struct{ myError.MyError }
	SetFormBodyError     struct{ myError.MyError }
	SetXmlBodyError      struct{ myError.MyError }
	SetJsonBodyError     struct{ myError.MyError }
	WriteResponseError   struct{ myError.MyError }
)

var (
	ReadResponseErr    ReadResponseError
	UrlEmptyErr        UrlEmptyError
	GenerateCertErr    GenerateCertError
	GenerateRequestErr GenerateRequestError
	UnmarshalXmlErr    UnmarshalXmlError
	UnmarshalJsonErr   UnmarshalJsonError
	SetSteamBodyErr    SetSteamBodyError
	SetFormBodyErr     SetFormBodyError
	SetXmlBodyErr      SetXmlBodyError
	SetJsonBodyErr     SetJsonBodyError
	WriteResponseErr   WriteResponseError
)

func (*ReadResponseError) New(msg string) myError.IMyError {
	return &ReadResponseError{MyError: myError.MyError{Msg: array.New([]string{"读取响应体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*ReadResponseError) Wrap(err error) myError.IMyError {
	return &ReadResponseError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "读取响应体失败", fmt.Errorf("读取响应体失败：%w", err).Error())}}
}

func (my *ReadResponseError) Error() string { return my.MyError.Msg }

func (my *ReadResponseError) Is(target error) bool {
	return reflect.DeepEqual(target, &ReadResponseErr)
}

func (*UrlEmptyError) New(msg string) myError.IMyError {
	return &UrlEmptyError{MyError: myError.MyError{Msg: array.New([]string{"url不能为空", msg}).JoinWithoutEmpty("：")}}
}

func (*UrlEmptyError) Wrap(err error) myError.IMyError {
	return &UrlEmptyError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "url不能为空", fmt.Errorf("url不能为空：%w", err).Error())}}
}

func (my *UrlEmptyError) Error() string { return my.MyError.Msg }

func (my *UrlEmptyError) Is(target error) bool {
	return reflect.DeepEqual(target, &UrlEmptyErr)
}

func (*GenerateCertError) New(msg string) myError.IMyError {
	return &GenerateCertError{MyError: myError.MyError{Msg: array.New([]string{"生成证书失败", msg}).JoinWithoutEmpty("：")}}
}

func (*GenerateCertError) Wrap(err error) myError.IMyError {
	return &GenerateCertError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "生成证书失败", fmt.Errorf("生成证书失败：%w", err).Error())}}
}

func (my *GenerateCertError) Error() string { return my.MyError.Msg }

func (my *GenerateCertError) Is(target error) bool {
	return reflect.DeepEqual(target, &GenerateCertErr)
}

func (*GenerateRequestError) New(msg string) myError.IMyError {
	return &GenerateRequestError{MyError: myError.MyError{Msg: array.New([]string{"生成请求对象失败", msg}).JoinWithoutEmpty("：")}}
}

func (*GenerateRequestError) Wrap(err error) myError.IMyError {
	return &GenerateRequestError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "生成请求对象失败", fmt.Errorf("生成请求对象失败：%w", err).Error())}}
}

func (my *GenerateRequestError) Error() string { return my.MyError.Msg }

func (my *GenerateRequestError) Is(target error) bool {
	return reflect.DeepEqual(target, &GenerateRequestErr)
}

func (*UnmarshalXmlError) New(msg string) myError.IMyError {
	return &UnmarshalXmlError{MyError: myError.MyError{Msg: array.New([]string{"获取xml格式响应体错误", msg}).JoinWithoutEmpty("：")}}
}

func (*UnmarshalXmlError) Wrap(err error) myError.IMyError {
	return &UnmarshalXmlError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "获取xml格式响应体错误", fmt.Errorf("获取xml格式响应体错误：%w", err).Error())}}
}

func (my *UnmarshalXmlError) Error() string { return my.MyError.Msg }

func (my *UnmarshalXmlError) Is(target error) bool {
	return reflect.DeepEqual(target, &UnmarshalXmlErr)
}

func (*UnmarshalJsonError) New(msg string) myError.IMyError {
	return &UnmarshalJsonError{MyError: myError.MyError{Msg: array.New([]string{"获取json格式响应体错误", msg}).JoinWithoutEmpty("：")}}
}

func (*UnmarshalJsonError) Wrap(err error) myError.IMyError {
	return &UnmarshalJsonError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "获取json格式响应体错误", fmt.Errorf("获取json格式响应体错误：%w", err).Error())}}
}

func (my *UnmarshalJsonError) Error() string { return my.MyError.Msg }

func (my *UnmarshalJsonError) Is(target error) bool {
	return reflect.DeepEqual(target, &UnmarshalJsonErr)
}

func (*SetSteamBodyError) New(msg string) myError.IMyError {
	return &SetSteamBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置二进制请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetSteamBodyError) Wrap(err error) myError.IMyError {
	return &SetSteamBodyError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "设置二进制请求体失败", fmt.Errorf("设置二进制请求体失败：%w", err).Error())}}
}

func (my *SetSteamBodyError) Error() string { return my.MyError.Msg }

func (my *SetSteamBodyError) Is(target error) bool {
	return reflect.DeepEqual(target, &SetSteamBodyErr)
}

func (*SetFormBodyError) New(msg string) myError.IMyError {
	return &SetFormBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置表单数据请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetFormBodyError) Wrap(err error) myError.IMyError {
	return &SetFormBodyError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "设置表单数据请求体失败", fmt.Errorf("设置表单数据请求体失败：%w", err).Error())}}
}

func (my *SetFormBodyError) Error() string { return my.MyError.Msg }

func (my *SetFormBodyError) Is(target error) bool {
	return reflect.DeepEqual(target, &SetFormBodyErr)
}

func (*SetXmlBodyError) New(msg string) myError.IMyError {
	return &SetXmlBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置xml请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetXmlBodyError) Wrap(err error) myError.IMyError {
	return &SetXmlBodyError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "设置xml请求体失败", fmt.Errorf("设置xml请求体失败：%w", err).Error())}}
}

func (my *SetXmlBodyError) Error() string { return my.MyError.Msg }

func (my *SetXmlBodyError) Is(target error) bool {
	return reflect.DeepEqual(target, &SetXmlBodyErr)
}

func (*SetJsonBodyError) New(msg string) myError.IMyError {
	return &SetJsonBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置json请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetJsonBodyError) Wrap(err error) myError.IMyError {
	return &SetJsonBodyError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "设置json请求体失败", fmt.Errorf("设置json请求体失败：%w", err).Error())}}
}

func (my *SetJsonBodyError) Error() string { return my.MyError.Msg }

func (my *SetJsonBodyError) Is(target error) bool {
	return reflect.DeepEqual(target, &SetJsonBodyErr)
}

func (*WriteResponseError) New(msg string) myError.IMyError {
	return &WriteResponseError{MyError: myError.MyError{Msg: array.New([]string{"写入响应失败", msg}).JoinWithoutEmpty("：")}}
}

func (*WriteResponseError) Wrap(err error) myError.IMyError {
	return &WriteResponseError{MyError: myError.MyError{Msg: operation.Ternary(err == nil, "写入响应失败", fmt.Errorf("写入响应失败：%w", err).Error())}}
}

func (my *WriteResponseError) Error() string { return my.MyError.Msg }

func (my *WriteResponseError) Is(target error) bool {
	return reflect.DeepEqual(target, &WriteResponseErr)
}
