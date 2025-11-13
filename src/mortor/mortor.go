package mortor

import (
	//"encoding/csv"
	//"encoding/json"

	"fmt"
	"sync"

	//"os"
	"net"

	//"github.com/valyala/fasthttp"
	"github.com/kafy11/gosocket/log"

	"motorv2/pkg/db"
	"motorv2/pkg/ws"
	awsC "motorv2/src/awsControllers"
	"motorv2/src/returnEng"
	"motorv2/src/sendSchedule"
)

type MotorConfigParams struct {
	WsBaseUrl  string
	WsAuthUser string
	WsAuthPass string
}

var motorConfig *MotorConfigParams

func SetConfig(params *MotorConfigParams) {
	motorConfig = params
}

// func Run(layoutStr, executionTimeStr string) error {
func Run() error {
	var macID string

	log.Info("Motor - Inicio")
	ws.SetBaseUrl(motorConfig.WsBaseUrl)
	ws.SetAuth(motorConfig.WsAuthUser, motorConfig.WsAuthPass)

	interfaces, err := net.Interfaces()
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

	dbInfos, err := awsC.GetInfos("Adm_get_wes_env?id=" + macID)
	if err != nil {
		fmt.Println("Erro de puxada no Controller Adm_get_wes_env:", err)
		log.Error("Erro de puxada no Controller Adm_get_wes_env:", err)
		return err
	}

	executionTimeStr := dbInfos["EXECUTION_TIME"].(string)
	fmt.Println("MOTOR EXECUTION_TIME:", dbInfos["EXECUTION_TIME"].(string))

	db.SetConfig(&db.DBConnectionParams{
		DB:     dbInfos["DB"].(string),
		DBUser: dbInfos["DB_USER"].(string),
		DBPass: dbInfos["DB_PASSWORD"].(string),
		DBHost: dbInfos["DB_HOST"].(string),
		DBPort: dbInfos["DB_PORT"].(string),
		DBSid:  dbInfos["DB_SID"].(string),
		DBName: dbInfos["DB_NAME"].(string),
	})
	log.Info("Banco configurado foi: " + dbInfos["DB"].(string))

	var wg sync.WaitGroup
	// Adiciona duas goroutines ao WaitGroup
	wg.Add(2)

	// Executa returnEng em uma goroutine
	go func() {
		defer wg.Done() // Marca como concluído ao final
		returnEng.ReturnAction()
	}()

	go func() {
		defer wg.Done() // Marca como concluído ao final
		sendSchedule.BusSendSchedule(executionTimeStr, macID)
	}()

	// Aguarda a conclusão de ambas as goroutines
	wg.Wait()

	fmt.Println("Todos os processos foram concluídas")
	return nil
}
