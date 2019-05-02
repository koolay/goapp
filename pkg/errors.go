// 自定义异常

package pkg

import (
	"fmt"

	"gopkg.in/go-playground/validator.v8"
)

type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func (err CommonError) Error() string {
	for field, msg := range err.Errors {
		return fmt.Sprintf("%s: %+v", field, msg)
	}
	return ""
}

// ErrInvalidForm 无效的表单输入
// 包括json, formdata, query
type ErrInvalidForm string

func (err ErrInvalidForm) Error() string {
	return string(err)
}

func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param != "" {
			res.Errors[v.Field] = fmt.Sprintf("{%v: %v}", v.Tag, v.Param)
		} else {
			res.Errors[v.Field] = fmt.Sprintf("{key: %v}", v.Tag)
		}

	}
	return res
}
