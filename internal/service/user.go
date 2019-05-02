package service

import (
	"context"

	"github.com/koolay/goapp/internal/store/entity"
)

var (
	defaultUserFields = "id, display_name, email, avatar, mobile, created_at, last_login_time"
)

type User struct {
	BaseServcie
}

func NewUser(base BaseServcie) *User {
	return &User{
		BaseServcie: base,
	}
}

func (s *User) GetProfile(ctx context.Context, userID string) (*entity.User, error) {
	return s.appCtx.Storage.GetUserByID(ctx, defaultUserFields, userID)
}
