package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"github.com/koolay/goapp/pkg"
)

func sendFailure(c *gin.Context, code int, errmsg string) {
	c.JSON(200, gin.H{"code": code, "message": errmsg})
}

func sendSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(200, gin.H{"code": code, "message": "ok", "data": data})
}

// bind 绑定表单，并验证
func bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

func checkError(c *gin.Context, err error) bool {
	if err == nil {
		return true
	}
	switch err.(type) {
	case *pkg.ErrInvalidForm:
		c.JSON(200, gin.H{"code": 400, "message": err.Error()})
		return false
	case validator.ValidationErrors:
		c.JSON(200, gin.H{"code": 400, "message": pkg.NewValidatorError(err).Error()})
		return false
	default:
		c.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return false
	}
}
