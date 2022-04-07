package hao

import (
	"fmt"
)

type authorizationError struct {
	operation string
	err       error
}

func (e *authorizationError) Error() string {
	return fmt.Sprintf("authorization failed during %s: %v", e.operation, e.err)
}

type causer interface {
	Cause() error
}

func (e *authorizationError) Cause() error {
	return e.err
}

func RunError() {
	// var err error
	// if err != nil {
	// 	return errors.Wrap(err, "read failed")
	// }

	// Cause 接口
	// switch err := errors.Cause(err).(type) {
	// case *MyError:
	// 	//
	// default:
	// 	//
	// }
}
