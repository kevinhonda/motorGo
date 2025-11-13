package db

import (
	"errors"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type MsSQLEngine struct {
	params *MsSQLConnectionParams
}

type MsSQLConnectionParams struct {
	User     string
	Password string
	Host     string
	Port     int
}

func NewMsSQLEngine(params *MsSQLConnectionParams) (*MsSQLEngine, error) {
	if params.Host == "" {
		return nil, errors.New("host cannot be empty")
	}

	if params.User == "" {
		return nil, errors.New("user cannot be empty")
	}

	if params.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if params.Port == 0 {
		return nil, errors.New("port cannot be empty")
	}

	connParams := &MsSQLConnectionParams{
		Host:     params.Host,
		User:     params.User,
		Password: params.Password,
		Port:     params.Port,
	}

	return &MsSQLEngine{connParams}, nil
}

func (o *MsSQLEngine) Connect() (*DB, error) {
	sqlxDB, err := sqlx.Connect("mssql",
		fmt.Sprintf(`server=%s;user id=%s;password=%s;port=%d`,
			o.params.Host,
			o.params.User,
			o.params.Password,
			o.params.Port),
	)
	if err != nil {
		return nil, err
	}

	return &DB{sqlxDB}, nil
}
