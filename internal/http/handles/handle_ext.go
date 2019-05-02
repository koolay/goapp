package handles

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"

	"github.com/koolay/goapp/internal/http/middleware"
	"github.com/koolay/goapp/internal/http/xerrors"
	"github.com/koolay/goapp/pkg"
)

func (h *Handler) GetSession(c *gin.Context) *middleware.Session {
	sess, exist := c.Get("sess")
	if exist {
		return sess.(*middleware.Session)
	}
	c.AbortWithStatus(401)
	return nil
}

func (h *Handler) SendFailure(c *gin.Context, code int, errmsg string) {
	c.JSON(200, gin.H{"code": code, "message": errmsg})
}

func (h *Handler) SendSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(200, gin.H{"code": code, "message": "ok", "data": data})
}

// MustBind 绑定表单，并验证
func (h *Handler) MustBind(c *gin.Context, obj interface{}) bool {
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.ShouldBindWith(obj, b)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": err.Error()})
		return false
	}
	return true
}

func (h *Handler) CheckError(c *gin.Context, err error) bool {
	if err == nil {
		return true
	}
	switch ue := err.(type) {
	case *xerrors.UserError:
		c.JSON(200, gin.H{
			"code":    ue.Code,
			"message": ue.Message,
		})
		return false
	case *pkg.ErrInvalidForm:
		c.JSON(200, gin.H{"code": 400, "message": err.Error()})
		return false
	case validator.ValidationErrors:
		c.JSON(200, gin.H{"code": 400, "message": pkg.NewValidatorError(err).Error()})
		return false
	default:
		h.logger.WithFields(logrus.Fields{
			"URL": c.Request.RequestURI,
		}).Error(fmt.Sprintf("%+v", err))
		c.JSON(500, "Server is busy.")
		return false
	}
}
