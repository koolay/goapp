// Package mysql provides database access
package mysql

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	Driver        string
	ConnectionURL string

	MaxIdleConns int
	MaxOpenConns int
}

// SQLDatabase mysql implement
type SQLDatabase struct {
	session       *sqlx.DB
	logger        *logrus.Logger
	configuration *Configuration
}

func TestConfiguration() Configuration {
	return Configuration{
		Driver:        "mysql",
		ConnectionURL: "root:dev@(127.0.0.1:3306)/goapp",
		MaxIdleConns:  0,
		MaxOpenConns:  4,
	}
}

func DefaultConfiguration() Configuration {
	return Configuration{
		Driver:        "mysql",
		ConnectionURL: "root@(localhost)/dev",
		MaxIdleConns:  0,
		MaxOpenConns:  4,
	}
}

// NewSQLDatabase init mysql Database instance
func NewSQLDatabase(logger *logrus.Logger, configuration Configuration) (*SQLDatabase, error) {
	var err error

	if mergo.Merge(&configuration, DefaultConfiguration()); err != nil {
		return nil, errors.Wrap(err, "Failed to merge configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug("Connecting db")
	var conn *sqlx.DB
	connCh := make(chan error)

	go func() {
		conn, err = sqlx.ConnectContext(ctx, configuration.Driver, configuration.ConnectionURL)
		if err != nil {
			err = errors.Wrapf(err, "Failed to ping: %s", configuration.ConnectionURL)
		}
		connCh <- err
	}()

	select {
	case <-time.After(5 * time.Second):
		return nil, errors.New("Failed to connect with timeout")
	case <-connCh:
		if err != nil {
			return nil, err
		}
	}

	conn.SetMaxOpenConns(configuration.MaxOpenConns)
	conn.SetMaxIdleConns(configuration.MaxIdleConns)

	mysqlDB := &SQLDatabase{
		configuration: &configuration,
	}
	mysqlDB.session = conn
	return mysqlDB, err
}
func (d *SQLDatabase) Close() error {
	if d.session != nil {
		return d.session.Close()
	}
	return nil
}

func (d *SQLDatabase) GetSession() *sqlx.DB {
	return d.session
}
