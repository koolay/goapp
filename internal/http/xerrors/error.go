package xerrors

import (
	"fmt"
)

// UserError 用户业务错误
type UserError struct {
	Code    int
	Message string
}

func (e *UserError) Error() string {
	return fmt.Sprintf("[UserError] code:%d, message:%s", e.Code, e.Message)
}
