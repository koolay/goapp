package store

import (
	"github.com/koolay/goapp/pkg/db/mysql"
)

type Storage struct {
	sqldb *mysql.SQLDatabase
}

func NewStorage(db *mysql.SQLDatabase) *Storage {
	return &Storage{
		sqldb: db,
	}
}
