package handles

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/koolay/goapp/internal/service/input"
)

// Register register
func (au *Handler) Register(c *gin.Context) {
	form := &input.Register{}
	if !au.MustBind(c, form) {
		return
	}

	userid, err := au.authService.Register(c, form)
	if !au.CheckError(c, err) {
		return
	}
	user, err := au.userService.GetProfile(c, userid)
	if !au.CheckError(c, err) {
		return
	}

	au.SendSuccess(c, 200, user)
}

// Login login
func (au *Handler) Login(c *gin.Context) {
	fmLogin := &input.Login{}
	if !au.MustBind(c, fmLogin) {
		return
	}
	user, err := au.authService.Login(c, fmLogin)
	if !au.CheckError(c, err) {
		return
	}

	exp := time.Now().Add(time.Second * time.Duration(au.appCfg.Security.JWTTimeout)).
		Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":          user.ID,
		"display_name": user.DisplayName,
		"exp":          exp,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(au.appCfg.Security.JWTSecret))
	if !au.CheckError(c, err) {
		return
	}

	c.SetCookie(au.appCfg.Security.TokenName, tokenString, 0, "/", "", false, true)
	au.SendSuccess(c, 200, gin.H{
		"token": tokenString,
		"user":  user,
	})

}

func (au *Handler) Logout(c *gin.Context) {
	c.SetCookie(au.appCfg.Security.TokenName, "", -1, "/", "", false, true)
	au.SendSuccess(c, 200, nil)
}
