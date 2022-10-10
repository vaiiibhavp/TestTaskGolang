package gopg

import (
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	Driver, Host, User, Pass, Name string
}

func NewSqlDB(dbConf *DbConfig) (*pg.DB, error) {

	switch dbConf.Driver {
	case POSTGRES:
		postgresqlDbInstance, err := NewPostgresqlDB(dbConf.Host, dbConf.User, dbConf.Pass, dbConf.Name)
		return postgresqlDbInstance, err
	default:
		dbConErr := NewDBConnError(UNSUPPORTED_DB_DRIVER)
		return nil, dbConErr
	}
}
