package appError

//自定义app错误

import (
	"fmt"
)

type Error struct {
	Code int32
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %v msg: %s ", e.Code, e.Msg)
}

func (e *Error) FormatMessage(values ...string) *Error {
	err := *e //clone 一份数据，防止修改原数据
	err.Msg = fmt.Sprintf(e.Msg, values)
	return &err
}

func NewError(message string) *Error {
	return &Error{Code: ServerError.Code, Msg: message}
}
