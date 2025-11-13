package db

import (
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresEngine struct {
	params *PostgresConnectionParams
}

type PostgresConnectionParams struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
}

func NewPostgresEngine(params *PostgresConnectionParams) (*PostgresEngine, error) {
	if params.Host == "" {
		return nil, errors.New("host cannot be empty")
	}

	if params.User == "" {
		return nil, errors.New("user cannot be empty")
	}

	if params.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if params.DBName == "" {
		return nil, errors.New("database name cannot be empty")
	}

	connParams := &PostgresConnectionParams{
		Host:     params.Host,
		User:     params.User,
		DBName:   params.DBName,
		Password: params.Password,
	}

	if params.Port != 0 {
		connParams.Port = params.Port
	} else {
		connParams.Port = 5432
	}

	return &PostgresEngine{connParams}, nil
}

func (o *PostgresEngine) Connect() (*DB, error) {
	sqlxDB, err := sqlx.Connect("pgx",
		fmt.Sprintf(`user=%s password=%s dbname=%s host=%s port=%d`,
			o.params.User,
			o.params.Password,
			o.params.DBName,
			o.params.Host,
			o.params.Port,
		),
	)
	if err != nil {
		return nil, err
	}

	return &DB{sqlxDB}, nil
}
