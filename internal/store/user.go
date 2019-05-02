package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/koolay/goapp/internal/store/entity"
)

// InsertUser insert user
func (s *Storage) InsertUser(ctx context.Context, user *entity.User) error {

	sqlTxt := `INSERT INTO users(id, account, password, mobile, display_name, email)
	VALUES (:id, :account, :password, :mobile, :display_name, :email)`

	_, err := s.sqldb.GetSession().NamedExecContext(ctx, sqlTxt, user)
	return err
}

// GetUserByAccount get user by account
func (s *Storage) GetUserByAccount(ctx context.Context, fields, account string) (*entity.User, error) {
	user := &entity.User{}
	if fields == "" {
		fields = "*"
	}
	var err error
	if err = s.sqldb.GetSession().GetContext(ctx, user, fmt.Sprintf(`SELECT %s FROM users WHERE account=?`, fields), account); err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByEmail get by email
func (s *Storage) GetUserByEmail(ctx context.Context, fields, email string) (*entity.User, error) {
	user := &entity.User{}
	if fields == "" {
		fields = "*"
	}
	var err error
	sqlTxt := fmt.Sprintf(`SELECT %s FROM users WHERE email=?`, fields)
	if err = s.sqldb.GetSession().GetContext(ctx, user, sqlTxt, email); err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByMobile get by mobile
func (s *Storage) GetUserByMobile(ctx context.Context, fields, mobile string) (*entity.User, error) {
	user := &entity.User{}
	if fields == "" {
		fields = "*"
	}
	var err error
	sqlTxt := fmt.Sprintf(`SELECT %s FROM users WHERE mobile=?`, fields)
	if err = s.sqldb.GetSession().GetContext(ctx, user, sqlTxt, mobile); err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByID get by id
func (s *Storage) GetUserByID(ctx context.Context, fields, userID string) (*entity.User, error) {
	user := &entity.User{}
	if fields == "" {
		fields = "*"
	}
	var err error
	sqlTxt := fmt.Sprintf(`SELECT %s FROM users WHERE id=?`, fields)
	if err = s.sqldb.GetSession().GetContext(ctx, user, sqlTxt, userID); err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// ExistAccount exist account if
func (s *Storage) ExistAccount(ctx context.Context, account string) (bool, error) {
	var userid string
	err := s.sqldb.GetSession().GetContext(ctx, &userid, "SELECT id FROM users WHERE account=?", account)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}

// ExistDisplayname exist displayName
func (s *Storage) ExistDisplayname(ctx context.Context, displayName string) (bool, error) {
	user := &entity.User{}
	err := s.sqldb.GetSession().GetContext(ctx, user, `select id from users where display_name=? limit 1`, displayName)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, err
}

// ExistMobile exist displayName
func (s *Storage) ExistMobile(ctx context.Context, mobile string) (bool, error) {
	user := &entity.User{}
	err := s.sqldb.GetSession().GetContext(ctx, user, `select id from users where mobile=? limit 1`, mobile)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}
