package mysql

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestNewMysqlDatabase(t *testing.T) {
	configuration := TestConfiguration()
	var err error

	db, err := NewSQLDatabase(logrus.New(), configuration)
	if err == nil {
		var dest string
		err = db.GetSession().Get(&dest, "select name from project limit 1")
		assert.Nil(t, err)
		assert.NotEmpty(t, dest)
	} else {
		t.Error(err)
	}
}
