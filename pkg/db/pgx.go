package db

import (
	"database/sql"
	"fmt"

	"github.com/bibaroc/wingman/pkg"
	"github.com/google/wire"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	PGXSet = wire.NewSet(
		pkg.DBConfigurationFromEnv,
		NewPGXConnector,
	)
)

type PGXConnector struct{}

func (p PGXConnector) Connect(dbConfig pkg.DBConfiguration) (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName,
	)
	database, err := sql.Open("pgx", dataSource)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func NewPGXConnector() PGXConnector {
	return PGXConnector{}
}
