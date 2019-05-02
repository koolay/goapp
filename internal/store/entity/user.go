package entity

import (
	"time"

	"github.com/koolay/goapp/pkg/types"
)

// User  user
type User struct {
	ID            string         `json:"id,omitempty" db:"id"`
	Account       string         `json:"account,omitempty" db:"account"`
	Password      string         `json:"password,omitempty" db:"password"`
	DisplayName   string         `json:"display_name,omitempty" db:"display_name"`
	Mobile        string         `json:"mobile,omitempty" db:"mobile"`
	Email         string         `json:"email,omitempty" db:"email"`
	LastLoginTime types.NullTime `json:"last_login_time,omitempty" db:"last_login_time"`
	Avatar        string         `json:"avatar,omitempty" db:"avatar"`
	Disabled      byte           `json:"disabled,omitempty" db:"disabled"`
	CreatedAt     time.Time      `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at,omitempty" db:"updated_at"`
}
