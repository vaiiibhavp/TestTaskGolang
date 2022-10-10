package gopg

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

func NewPostgresqlDB(dbHost string, dbUser string, dbPass string, dbName string) (*pg.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, POSTGRESPORT, dbName)
	opt, err := pg.ParseURL(connectionString)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)
	return db, nil
}
