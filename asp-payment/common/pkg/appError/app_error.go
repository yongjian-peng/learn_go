package appError

//自定义app错误

import (
	"fmt"
)

type Error struct {
	Code    int32
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %v msg: %s ", e.Code, e.Message)
}

func (e *Error) FormatMessage(values ...interface{}) *Error {
	err := *e //clone 一份数据，防止修改原数据
	err.Message = fmt.Sprintf(e.Message, values...)
	return &err
}

func NewError(message string) *Error {
	return &Error{Code: CodeUnknown.Code, Message: message}
}
