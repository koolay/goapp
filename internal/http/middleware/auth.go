package middleware

import (
	"fmt"

	"github.com/koolay/goapp/conf"
	"github.com/sirupsen/logrus"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Session struct {
	UserID string
}

type Authorization struct {
	AppCfg *conf.AppConfiguration
	Logger *logrus.Logger
}

func (au Authorization) AbortWithNoAuthorizated(c *gin.Context) {
	c.SetCookie(au.AppCfg.Security.TokenName, "", -1, "/", "", false, true)
	c.String(401, "UnAuthorized")
	c.Abort()
}

// AuthRequired require login
func (au Authorization) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		tokenString, err := c.Cookie(au.AppCfg.Security.TokenName)
		if err != nil {
			au.AbortWithNoAuthorizated(c)
			return
		}
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(au.AppCfg.Security.JWTSecret), nil
		})
		if err != nil {
			au.Logger.WithError(err).Error("Failed to parese jwt")
			au.AbortWithNoAuthorizated(c)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sess := &Session{
				UserID: claims["uid"].(string),
			}
			c.Set("sess", sess)
			c.Next()
			return
		}
		au.AbortWithNoAuthorizated(c)
	}
}
