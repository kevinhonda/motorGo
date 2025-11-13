package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/kafy11/gosocket/log"

	/*"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/db"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/ws"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/build"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/companyClient"*/
	"motorv2/pkg/ws"
	awsC "motorv2/src/awsControllers"
	"motorv2/src/websocket/pkg/db"
	"motorv2/src/websocket/pkg/wsWebsocket"
	"motorv2/src/websocket/src/build"
	"motorv2/src/websocket/src/companyClient"
)

var macID string
var params map[string]interface{}

func init() {
	var interfaces []net.Interface
	var err error

	interfaces, err = net.Interfaces()
	if err != nil {
		fmt.Println("Erro ao buscar interfaces de rede:", err)
		log.Error("Erro ao buscar interfaces de rede:", err)
	}

	// Loop pelas interfaces para encontrar o endereço MAC
	for _, iface := range interfaces {
		// Ignora interfaces loopback e interfaces sem endereço MAC
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			macAddr := iface.HardwareAddr
			fmt.Println("Endereço MAC:", macAddr.String())
			log.Info("Endereço MAC:" + macAddr.String())
			macID = macAddr.String()
		}
	}
}

func getParams() {

	ws.SetBaseUrl(os.Getenv("WS_BASE_URL"))
	ws.SetAuth(os.Getenv("WS_AUTH_USER"), os.Getenv("WS_AUTH_PASS"))
	log.Info("Puxando os parâmetros do websocket")
	dbInfos, err := awsC.GetInfos("Adm_get_wes_env?id=" + macID)
	if err != nil {
		fmt.Println("Erro de puxada no Controller Adm_get_wes_env para o websocket:", err)
		log.Error("Erro de puxada no Controller Adm_get_wes_env para o websocket:", err)
		return
	}
	params = dbInfos
	log.Info("Puxada no Controller Adm_get_wes_env para o websocket com sucesso")
}

func GetWSClientParams() *companyClient.Params {

	getParams()

	log.Info(params)

	return &companyClient.Params{
		URL:             params["WEBSOCKET_SERVER_URL"].(string),
		CompanyName:     params["COMPANY_NAME"].(string),
		DB:              params["DB"].(string),
		AuthUser:        params["WS_AUTH_USER"].(string),
		AuthPass:        params["WS_AUTH_PASS"].(string),
		Version:         build.Version,
		ControllersPath: params["CONTROLLERS_DIR_PATH"].(string),
	}
}

func ConnectDB() error {
	port, err := strconv.Atoi(params["DB_PORT"].(string))
	if err != nil {
		port = 0
	}

	var engine db.DBEngine
	dbType := strings.ToLower(params["DB"].(string))
	user := params["DB_USER"].(string)
	password := params["DB_PASSWORD"].(string)
	host := params["DB_HOST"].(string)

	switch dbType {
	case "oracle":
		log.Info("Conectando no banco Oracle")
		//connectionString := params["DB_CONNECTION_STRING"].(string)
		engine, err = db.NewOracleEngine(&db.OracleConnectionParams{
			User:             user,
			Password:         password,
			Host:             host,
			Port:             port,
			Sid:              params["DB_SID"].(string),
			ConnectionString: "", //connectionString,
		})
	case "postgres":
		engine, err = db.NewPostgresEngine(&db.PostgresConnectionParams{
			User:     user,
			Password: password,
			DBName:   params["DB_NAME"].(string),
			Host:     host,
			Port:     port,
		})
	case "sqlserver":
		engine, err = db.NewMsSQLEngine(&db.MsSQLConnectionParams{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
		})
	default:
		log.Error("Erro no .env: Tipo de banco inválido ou não informado")
		return errors.New("erro no .env: Tipo de banco inválido ou não informado")
	}

	if err != nil {
		log.Error("Falha ao criar engine do banco", err)
		return errors.New(fmt.Sprint("Falha ao criar engine do banco", err))

	}

	log.Info("Conectando no banco")
	if err := db.Connect(engine); err != nil {
		log.Error("Falha ao conectar no banco", err)
		return errors.New(fmt.Sprint("Falha ao conectar no banco", err))
	}
	log.Info("Conectado no banco com sucesso")

	return nil
}

func SetWebserviceParams() {
	wsWebsocket.SetAuth(os.Getenv("WS_AUTH_USER"), os.Getenv("WS_AUTH_PASS"))
	wsWebsocket.SetBaseUrl(os.Getenv("WS_BASE_URL"))
}
