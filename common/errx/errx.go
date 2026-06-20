package errx

import "fmt"

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func NewError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

var (
	ErrSuccess          = NewError(0, "success")
	ErrInvalidParam     = NewError(10001, "invalid parameter")
	ErrUnauthorized     = NewError(10002, "unauthorized")
	ErrInternalServer   = NewError(10003, "internal server error")
	ErrUserNotFound     = NewError(20001, "user not found")
	ErrUserExists       = NewError(20002, "user already exists")
	ErrPasswordWrong    = NewError(20003, "password is wrong")
	ErrProductNotFound  = NewError(30001, "product not found")
	ErrStockNotEnough   = NewError(40001, "stock not enough")
	ErrStockNotFound    = NewError(40002, "stock not found")
	ErrOrderNotFound    = NewError(50001, "order not found")
	ErrOrderCreateFail  = NewError(50002, "order create fail")
	ErrLockTimeout      = NewError(60001, "lock timeout")
	ErrKafkaSendFail    = NewError(70001, "kafka send fail")
)
