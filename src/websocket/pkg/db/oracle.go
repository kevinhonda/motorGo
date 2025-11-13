package db

import (
	"errors"
	"fmt"

	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"
)

type OracleEngine struct {
	params *OracleConnectionParams
}

type OracleConnectionParams struct {
	User             string
	Password         string
	Host             string
	Port             int
	Sid              string
	ConnectionString string
}

func NewOracleEngine(params *OracleConnectionParams) (*OracleEngine, error) {
	if params.Host == "" && params.ConnectionString == "" {
		return nil, errors.New("please inform an db host or connection string")
	}

	connParams := &OracleConnectionParams{
		Host:             params.Host,
		ConnectionString: params.ConnectionString,
	}

	if params.User != "" {
		connParams.User = params.User
	} else {
		connParams.User = "root"
	}

	if params.Password != "" {
		connParams.Password = params.Password
	} else {
		connParams.Password = "root"
	}

	if params.Port != 0 {
		connParams.Port = params.Port
	} else {
		connParams.Port = 1521
	}

	if params.Sid != "" {
		connParams.Sid = params.Sid
	} else {
		connParams.Sid = "ORCL"
	}
	return &OracleEngine{connParams}, nil
}

func (o *OracleEngine) Connect() (*DB, error) {
	var dataSourceName string

	if o.params.ConnectionString != "" {
		dataSourceName = fmt.Sprintf(`user="%s" password="%s" connectString="%s"`,
			o.params.User,
			o.params.Password,
			o.params.ConnectionString,
		)
	} else {
		dataSourceName = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`,
			o.params.User,
			o.params.Password,
			o.params.Host,
			o.params.Port,
			o.params.Sid,
		)
	}

	log.Info("WBsocket tentando conectar ao oracle: ", dataSourceName)

	sqlxDB, err := sqlx.Connect("godror", dataSourceName)
	if err != nil {
		log.Error("WBsocket error: ", err)
		return nil, err
	}

	return &DB{sqlxDB}, nil
}
