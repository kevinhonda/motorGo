package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"
	_ "github.com/lib/pq"
	_ "github.com/sijms/go-ora/v2"
)

type DBConnectionParams struct {
	DataSourceName string
	Connection     string
	DB             string
	DBUser         string
	DBPass         string
	DBHost         string
	DBPort         string
	DBSid          string
	DBName         string
}

var db *sqlx.DB

//var tx *sqlx.Tx

var dbConnection *DBConnectionParams

func SetConfig(params *DBConnectionParams) {
	dbConnection = params
}

func ConnectDB(dbName string) (*DBConnectionParams, error) {

	var dataSourceName string
	var dbConnect string
	/*
		log.Info("Engine - User:", dbConnection.DBUser)
		log.Info("Engine - Pass:", dbConnection.DBPass)
		log.Info("Engine - Host:", dbConnection.DBHost)
		log.Info("Engine - Port:", dbConnection.DBPort)
		log.Info("Engine - SID:", dbConnection.DBSid)
		log.Info("Engine - Name:", dbConnection.DBName)
	*/
	switch dbName {
	case "SQLserver":
		//dataSourceName = "server=" + dbConnection.DBHost + ";user id=" + dbConnection.DBUser + ";password=" + dbConnection.DBPass + ";database=" + dbConnection.DBName
		dataSourceName = "server=" + dbConnection.DBHost + ";user id=" + dbConnection.DBUser + ";password=" + dbConnection.DBPass + ";database=" + dbConnection.DBName
		dbConnect = "sqlserver"
		return &DBConnectionParams{
			DataSourceName: dataSourceName,
			Connection:     dbConnect,
		}, nil
	case "Postgre":
		//dataSourceName = "user=" + dbConnection.DBUser + "password=" + dbConnection.DBPass + "dbname=" + dbConnection.DBName + "host=" + dbConnection.DBHost + "port=" + dbConnection.DBPort
		dataSourceName = "host=localhost user=" + dbConnection.DBUser + " dbname=" + dbConnection.DBName + " sslmode=disable password=" + dbConnection.DBPass

		dbConnect = "postgres"
		return &DBConnectionParams{
			DataSourceName: dataSourceName,
			Connection:     dbConnect,
		}, nil
	case "Oracle":
		encodedPassword := url.QueryEscape(dbConnection.DBPass)
		dataSourceName := "oracle://" + dbConnection.DBUser + ":" + encodedPassword + "@" + dbConnection.DBHost + ":" + dbConnection.DBPort + "/" + dbConnection.DBSid + ""

		log.Info("Engine - dataSourceName:", dataSourceName)
		dbConnect := "oracle"
		return &DBConnectionParams{
			DataSourceName: dataSourceName,
			Connection:     dbConnect,
		}, nil
	default:
		return nil, errors.New("Erro no .env: Tipo de banco inválido ou não informado")
		//return nil, nil
	}
	//return nil, nil
}

// func ConnCheck(dbName string, dbShits map[string]string) (string, error) {
func ConnCheck() (string, error) {
	dbConnectionParams, err := ConnectDB(dbConnection.DB)
	//log.Info("ConnCheck dbConnectionParams - DataSourceName: ", dbConnectionParams.DataSourceName)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		log.Error("DataSourceName: ", dbConnectionParams.DataSourceName)
		log.Error("Erro ao conectar ao banco de dados: ", err)
		return "Erro ao conectar ao banco de dados:", err
	}

	// Abre uma conexão
	log.Info("sql.Open - Connection: ", dbConnectionParams.Connection)
	log.Info("sql.Open - DataSourceName: ", dbConnectionParams.DataSourceName)
	db, err := sql.Open(dbConnectionParams.Connection, dbConnectionParams.DataSourceName)
	if err != nil {
		log.Error("DataSourceName: ", dbConnectionParams.DataSourceName)
		log.Error("Erro ao abrir a conexão: ", err)
		return "Erro ao abrir a conexão:", err
	}
	defer db.Close()

	////////////////////////////////////////////////////////////////////
	// Testa a conexão
	err = db.Ping()
	if err != nil {
		log.Error("DataSourceName: ", dbConnectionParams.DataSourceName)
		log.Error("Erro ao testar a conexão: ", err)
		return "Erro ao testar a conexão:", err
	}
	return "Conexão com o banco bem-sucedida!", nil
}

func OpenDB() (*sqlx.DB, error) {

	log.Info("Banco: " + dbConnection.DB)
	dbConnection, err := ConnectDB(dbConnection.DB)
	if err != nil {
		log.Error("Erro de conexão de banco: ", err)
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		return nil, err
	}

	//db, err = sqlx.Open(dbConnection.Connection, dbConnection.DataSourceName+"connect timeout=300&"+"read timeout=1800&"+"write timeout=1800")
	db, err = sqlx.Open(dbConnection.Connection, dbConnection.DataSourceName)
	if err != nil {
		log.Error("Erro ao conectar ao banco de dados:", err)
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		//return  err
	}

	// Check if the connection is open by pinging the database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func OpenTX() (*sqlx.Tx, error) {
	tx, err := db.Beginx()
	if err != nil {
		log.Error("Erro ao iniciar transação: ", err)
		return nil, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	return tx, nil
}

// func SqlExec(query string, fileArray map[string]interface{}) (*sqlx.Rows, string) {
func SqlExec(tx *sqlx.Tx, query string, fileArray map[string]interface{}) (int, string) {
	fmt.Println("Executando query")
	//total, err := db.NamedExec(queryAction, fileArray)
	//fmt.Printf("O tipo da variável total é: %T\n", total)
	errorMessage := ""

	log.Info("SqlExec: " + dbConnection.DB)
	if dbConnection.DB == "Oracle" {
		finalQuery := replacePlaceholders(query, fileArray)

		log.Info("Query Final:")
		log.Info(finalQuery)

		var result int
		_, err := tx.Exec(finalQuery, sql.Named("result", &result))

		if err != nil {
			log.Error("SqlExec err: ", err)
			log.Error(finalQuery)
			errorMessage = err.Error()
			return 0, errorMessage
		}

		fmt.Println("ID gerado:", result)
		log.Info("SqlExec - ID gerado:", result)
		return result, errorMessage
	}
	rows, err := tx.NamedQuery(query, fileArray)
	log.Info("NamedQuery - Executando query:", query)
	log.Info("NamedQuery - Executando fileArray:", fileArray)
	if err != nil {
		//log.Fatal("Error executing query: ", err.Error())
		log.Error("Error executing query: ", err.Error())
		errorMessage = err.Error()
		//return 0, fmt.Sprintf("Erro: %s", err)
	}

	resultado := 0

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(&resultado)
			if err != nil {
				log.Fatal("Error scanning row: ", err.Error())
			}
		}
	}

	return resultado, errorMessage

	//return rows, errorMessage

}

func SqlRun(query, layout string) (*sqlx.Rows, error) {
	// Log da query (cuidado com dados sensíveis)
	log.Info(fmt.Sprintf("\n%s: \n %s\n", layout, query))

	// Configurar contexto
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)

	// Executar query
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		cancel() // Liberar recursos do contexto em caso de erro
		log.Error(layout, " - Erro ao executar a consulta:", err)
		erroLog := fmt.Sprintf("%s - Erro ao executar a consulta: %s \n QUERY: \n %s", layout, err, query)
		log.Error(erroLog)
		fmt.Println(layout, " - Erro ao executar a consulta:", err)
		return nil, err
	}

	// Retorna rows e indica que há resultados
	// O chamador precisará processar a primeira linha e continuar iteração
	return rows, nil
}

func getValueFromMap(data map[string]interface{}, key string) string {
	if val, exists := data[key]; exists {
		switch v := val.(type) {
		case string:
			if v == "" {
				return "'" // Retornar '' se o valor for uma string vazia
			}
			return v
		case int, int64, float64:
			return fmt.Sprintf("%v", v)
		case []interface{}:
			// Se for uma lista, pode ser necessário acessar o primeiro item, por exemplo.
			if len(v) > 0 {
				return fmt.Sprintf("%v", v[0])
			}
		default:
			return "" // Retornar vazio para tipos não esperados
		}
	}
	return "" // Retornar vazio se a chave não existir
}

// Função para substituir os placeholders da query pelos valores do JSON
func replacePlaceholders(query string, jsonData map[string]interface{}) string {
	// Regex para encontrar os placeholders no formato :param
	re := regexp.MustCompile(`:\w+`)

	// Função para substituir os placeholders
	replacedQuery := re.ReplaceAllStringFunc(query, func(placeholder string) string {
		// Remove o ":" do placeholder
		param := placeholder[1:]

		// Busca o valor correspondente no map
		if val := getValueFromMap(jsonData, param); val != "" {
			// Tratar números e strings
			if strings.HasPrefix(placeholder, "to_number(") {
				return val
			} else if val == "'" {
				return ""
			} else {
				return fmt.Sprintf("%s", val)
			}
		}

		// Se o parâmetro não estiver no JSON, retorna o placeholder original
		return placeholder
	})

	return replacedQuery
}
