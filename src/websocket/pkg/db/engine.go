package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"
)

type DB struct {
	*sqlx.DB
}

var Connection *DB

type DBEngine interface {
	Connect() (*DB, error)
}

func Connect(engine DBEngine) error {
	if conn, err := engine.Connect(); err != nil {
		log.Error("WBsocket ENGINE error: ", err)
		return err
	} else {
		log.Info("WBsocket ENGINE Conectado")
		Connection = conn
	}
	return nil
}

func (db *DB) MapSelect(query string, rowsLimit int, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Queryx(query, args...)
	if err != nil {
		log.Error("WBsocket ENGINE error: ", err)
		return nil, err
	}

	var result []map[string]interface{}
	count := 0
	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.MapScan(row)

		if err != nil {
			return nil, err
		}

		result = append(result, row)

		count++
		if count == rowsLimit {
			break
		}
	}
	return result, nil
}
