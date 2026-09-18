package errno

import (
	"errors"
	"fmt"
)

type Err struct {
	code    int64
	message string
	cause   error
}

func (e Err) Error() string {
	return fmt.Sprintf("[%d]%s", e.code, e.message)
}

func NewErr(code int64, message string) *Err {
	return &Err{
		code:    code,
		message: message,
	}
}

func (e Err) Code() int64 {
	return e.code
}

func (e Err) Message() string {
	return e.message
}

func (e Err) Unwrap() error {
	return e.cause
}

func Wrap(e *Err, cause error) *Err {
	return &Err{
		code:    e.code,
		message: e.message,
		cause:   cause,
	}
}

func WithMsg(e *Err, message string) Err {
	return Err{
		code:    e.code,
		message: message,
		cause:   e.cause,
	}
}

func Convert(err error) *Err {
	if err == nil {
		return nil
	}

	var bizErr *Err
	if errors.As(err, &bizErr) {
		return bizErr
	}

	return Wrap(UnknownError, err)
}
