package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/koolay/goapp/internal/http/xerrors"
	"github.com/koolay/goapp/internal/service/input"
	"github.com/koolay/goapp/internal/store/entity"
	"github.com/koolay/goapp/pkg"
)

type Auth struct {
	BaseServcie
}

func NewAuth(base BaseServcie) *Auth {
	return &Auth{
		BaseServcie: base,
	}
}

// Register register user
func (au *Auth) Register(ctx context.Context, registerInput *input.Register) (string, error) {
	var err error

	if registerInput.Password != registerInput.ConfirmPassword {
		return "", &xerrors.UserError{
			Code:    400,
			Message: "Password nt equals with confirmPassword",
		}
	}

	fmt.Println(pkg.PrimaryKey())
	var exist bool
	exist, err = au.appCtx.Storage.ExistAccount(ctx, registerInput.Account)
	if err != nil {
		return "", err
	}

	if exist {
		return "", &xerrors.UserError{
			Code:    400,
			Message: fmt.Sprintf("Account %s already exist", registerInput.Account),
		}
	}

	exist, err = au.appCtx.Storage.ExistDisplayname(ctx, registerInput.DisplayName)
	if err != nil {
		return "", err
	}

	if exist {
		return "", &xerrors.UserError{
			Code:    400,
			Message: fmt.Sprintf("display_name %s already exist", registerInput.DisplayName),
		}
	}

	user := &entity.User{
		ID:          pkg.PrimaryKey(),
		DisplayName: registerInput.DisplayName,
		Account:     registerInput.Account,
		Password:    registerInput.Password,
	}

	user.Password, err = pkg.EncryptPassword(user.Password)
	if err != nil {
		return "", errors.Wrapf(err, "Failed encrypt password: %s", user.Password)
	}

	return user.ID, au.appCtx.Storage.InsertUser(ctx, user)
}

// Login user
func (au *Auth) Login(ctx context.Context, loginInput *input.Login) (*entity.User, error) {
	if loginInput.Account == "" || loginInput.Password == "" {
		return nil, &xerrors.UserError{
			Code:    400,
			Message: "account or password is empty",
		}
	}

	user, err := au.appCtx.Storage.GetUserByAccount(ctx, "*", loginInput.Account)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &xerrors.UserError{
			Code:    404,
			Message: "Account and Password not match",
		}
	}

	if !pkg.ComparePassword(user.Password, loginInput.Password) {
		return nil, &xerrors.UserError{
			Code:    404,
			Message: "Account and Password not match",
		}
	}
	if user.Disabled == 1 {
		return nil, &xerrors.UserError{
			Code:    403,
			Message: "You had be forbid. Please contact the admin",
		}
	}
	user.Password = "*"
	return user, nil
}

// GetUserByAccount get user info
func (au *Auth) GetUserByAccount(ctx context.Context, account string) (*entity.User, error) {
	return au.appCtx.Storage.GetUserByAccount(ctx, "", account)
}
