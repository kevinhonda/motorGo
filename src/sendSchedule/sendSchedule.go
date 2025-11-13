package sendSchedule

import (
	"sync"
	//"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kafy11/gosocket/log"
	//"github.com/valyala/fasthttp"
	"github.com/jmoiron/sqlx"

	"motorv2/pkg/db"
	awsC "motorv2/src/awsControllers"
	"motorv2/src/sendStuffs"
)

var dbs *sqlx.DB

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func BusSendSchedule(executionTimes, macID string) {
	var executedTime string
	executionTimeStr := executionTimes
	fmt.Println("sendSchedule executionTimeStr:", executionTimeStr)
	log.Info("Exec Times: " + executionTimeStr)
	for {
		//fmt.Sprintf("Horario base de exec %s ", executionTimeStr)
		var durationUntilTarget time.Duration
		executeTimes := strings.Split(executionTimeStr, ",")

		for i := 0; i < len(executeTimes); i++ {
			executeTimes[i] = strings.ReplaceAll(executeTimes[i], " ", "")
			if executeTimes[i] == "" {
				log.Error("A variável de ambiente EXECUTION_TIME não está definida.")
				fmt.Println("A variável de ambiente EXECUTION_TIME não está definida.")
				//return errors.New("A variável de ambiente EXECUTION_TIME não está definida.")
			}

			//Divide a variável em hora e minutos
			timeParts := strings.Split(executeTimes[i], ":")

			if len(timeParts) != 2 {
				fmt.Println("timeParts")
				fmt.Println(timeParts)
				log.Error("A variável de ambiente EXECUTION_TIME não está no formato correto (HH:MM).")
				fmt.Println("A variável de ambiente EXECUTION_TIME não está no formato correto (HH:MM).")
				//return errors.New("A variável de ambiente EXECUTION_TIME não está no formato correto (HH:MM).")
			}
			//Analisa a hora e os minutos da variável de ambiente
			executionHour, err := strconv.Atoi(timeParts[0])

			if err != nil {
				log.Error(fmt.Sprint("Erro ao analisar a hora de execução:", err))
				fmt.Println(fmt.Sprint("Erro ao analisar a hora de execução:", err))
				//return errors.New(fmt.Sprint("Erro ao analisar a hora de execução:", err))
			}
			executionMinute, err := strconv.Atoi(timeParts[1])

			if err != nil {
				log.Error(fmt.Sprint("Erro ao analisar a minutos de execução:", err))
				fmt.Println(fmt.Sprint("Erro ao analisar a minutos de execução:", err))
				//return errors.New(fmt.Sprint("Erro ao analisar a minutos de execução:", err))
			}
			//Pega a data atual com o horario do parametro current e compara com o do parametro

			currentTime := time.Now()
			targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), executionHour, executionMinute, 0, 0, time.Local)
			log.Info("targetTime: ", targetTime)

			//Calcula a diferença entre o horario atual X horario do parametro
			//durationUntilTarget := targetTime.Sub(currentTime)
			newDurationUntilTarget := targetTime.Sub(currentTime)
			log.Info("newDurationUntilTarget: ", newDurationUntilTarget)

			//Verifica se a diferença é menor que zero, se for menor que zero quer dizer que o horario já passou
			//if durationUntilTarget < 0 {
			if newDurationUntilTarget < 0 {
				//Se a diferença(durationUntilTarget) for menor que zero ele add 24h no target(Data de execução)
				targetTime = targetTime.Add(24 * time.Hour)
				log.Info("New targetTime: ", targetTime)
				//Compara a proxima data de execução com a data atual
				//durationUntilTarget = targetTime.Sub(currentTime)
				newDurationUntilTarget = targetTime.Sub(currentTime)
				log.Info("New durationUntilTarget: ", newDurationUntilTarget)
			}
			if durationUntilTarget > newDurationUntilTarget || durationUntilTarget == 0 {
				durationUntilTarget = newDurationUntilTarget
				executedTime = executeTimes[i]
			}
			//Joga o durationUntilTarget(Diferença de tempo) para a execução
		}

		timer := time.NewTimer(durationUntilTarget)
		//log.Info("timer: ", timer)
		<-timer.C

		executionTimeStr = BusConfig(macID, "")
		fmt.Printf("Novo horario base de exec %s ", executionTimeStr)
		log.Info(fmt.Sprintf("Novo horario base de exec %s ", executionTimeStr))
		log.Info("Executado as " + executedTime)

	}
}

func BusConfig(macID, layoutVar string) string {
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Configurando paradas e chamando serviço que cria pre-signed e puxa query
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Capturar o tempo de início
	log.Info("////////////////////////////////////////////////////////////")
	log.Info("Envio de dados - Começo da execução. \n")
	log.Info("////////////////////////////////////////////////////////////")
	start := time.Now()

	dir, err := os.Getwd()
	// Verifica se ocorreu algum erro
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho:", err)
		log.Error("Erro ao obter o diretório de trabalho:", err)
		// Encerra o programa com código de saída 1 em caso de erro
		os.Exit(1)
	}

	// Imprime o diretório de trabalho atual
	fmt.Println("Diretório de trabalho atual:", dir)

	//Puxando configuracoes
	log.Info("Puxando configurações")
	if macID == "00:00:00:00:00:00:00:e0" {
		macID = ""
	}
	motorConfigs, err := awsC.GetInfos("Adm_get_wes_config?id=" + macID)
	if err != nil {
		log.Error(fmt.Sprint("GetInfos Error:", err))
		fmt.Println("BusConfig GetInfos Error:", err)
		return ""
	}
	log.Info("Setando querys")
	dataQuerys := motorConfigs["query"].(map[string]interface{})

	//Atribuindo valor para as variaveis
	dbInfos := motorConfigs["env"].(map[string]interface{})
	company := motorConfigs["COMPANY"].(string)

	//Atribuindo valor para as variaveis de query por company e por ERP
	// Verificar se a chave "DATA" existe dentro de "RESULT"
	if newQuery, exists := dataQuerys["NEW"].(map[string]interface{}); exists {
		var layoutStr string
		if layoutVar == "" {
			layoutStr = dbInfos["LAYOUTS"].(string)
		} else {
			layoutStr = layoutVar
		}

		layoutStr = strings.ReplaceAll(strings.ToUpper(layoutStr), " ", "")
		log.Info("Layouts: " + layoutStr)
		layouts := strings.Split(layoutStr, ",")

		//Configurando conexao com banco
		db.SetConfig(&db.DBConnectionParams{
			DB:     dbInfos["DB"].(string),
			DBUser: dbInfos["DB_USER"].(string),
			DBPass: dbInfos["DB_PASSWORD"].(string),
			DBHost: dbInfos["DB_HOST"].(string),
			DBPort: dbInfos["DB_PORT"].(string),
			DBSid:  dbInfos["DB_SID"].(string),
			DBName: dbInfos["DB_NAME"].(string),
		})

		fmt.Println("Banco configurado foi")
		log.Info("Banco configurado foi: " + dbInfos["DB"].(string))

		fmt.Println("SendStuff - New exec schedule " + dbInfos["EXECUTION_TIME"].(string))

		//LayoutSendGoutine(layouts, newQuery, company)
		//////////////////////////////////////////
		//GO LOOP TESTE
		//////////////////////////////////////////

		var wg sync.WaitGroup

		if dbs, err = db.OpenDB(); err != nil {
			fmt.Println("Error initializing database connection:", err)
		}

		for i := 0; i < len(layouts); i++ {
			wg.Add(1)
			go func(layout string) {
				defer wg.Done()

				fmt.Println(layout + ": Rodando o Layout mandando para BusSend")
				query, ok := newQuery[layout].(string)

				if ok && query != "" {
					query = strings.ReplaceAll(query, ":id_company", company)
					sendStuffs.BusSend(layout, query)
				} else {
					fmt.Println("O layout " + layout + " não possui query")
					log.Error("O layout " + layout + " não possui query")
				}
				fmt.Println("END: " + layout + " Antes layouts[i]")

				// Capturar o tempo de término
				layoutEnd := time.Now()
				layoutDuration := layoutEnd.Sub(start).Seconds()
				// Exibir a duração
				//fmt.Sprintf("Horario base de exec %s ", executionTimeStr)
				//fmt.Printf("O envio do layout %s levou %f segundos para executar.\n", layout, layoutDuration)
				log.Info(fmt.Sprintf("O envio do layout %s levou %f segundos para executar.\n", layout, layoutDuration))
			}(layouts[i])
		}

		// Espera todas as goroutines terminarem
		wg.Wait()

		defer dbs.Close()
		fmt.Println("DB Closed")
		log.Info("sendSchedule - DB Closed")
		/**/
		// Capturar o tempo de término
		end := time.Now()

		// Calcular a duração em segundos
		duration := end.Sub(start).Seconds()
		// Exibir a duração
		//fmt.Sprintf("Horario base de exec %s ", executionTimeStr)
		fmt.Printf("O envio levou %f segundos para executar.\n", duration)
		log.Info(fmt.Sprintf("O envio completo levou %f segundos para executar.\n", duration))

		log.Info("////////////////////////////////////////////////////////////")
		log.Info("Envio de dados - Fim da execução. \n")
		log.Info("////////////////////////////////////////////////////////////")
	} else {
		log.Error("Não possui querys de envio para a company ou para o layout.")
	}
	//Enviando Logs para S3
	/**/
	if fileExists("error_log.txt") {
		//fmt.Printf("O arquivo %s existe na pasta atual.\n", "log.error.txt")

		body, err := awsC.GetInfos("Adm_get_pre_signed?folder=erp_inventory&file_name=error_log")
		if err != nil {
			log.Error(fmt.Sprintf("Log send: GetInfos - %s ", err))
			fmt.Println("BusConfig GetInfos error_log - ", err)
		}

		if err := awsC.SendToS3("error_log.txt", body[strconv.Itoa(1)].(string)); err != nil {
			log.Error(fmt.Sprintf("Log send: SendToS3 - %s ", err))
			fmt.Println("BusConfig SendToS3 error_log.txt - ", err)
		}

		body, err = awsC.GetInfos("Adm_get_pre_signed?folder=erp_inventory&file_name=info_log")
		if err != nil {
			log.Error(fmt.Sprintf("Log send: GetInfos - %s ", err))
			fmt.Println("BusConfig GetInfos info_log - ", err)
		}

		if err := awsC.SendToS3("info_log.txt", body[strconv.Itoa(1)].(string)); err != nil {
			log.Error(fmt.Sprintf("Log send: SendToS3 - %s ", err))
			fmt.Println("BusConfig SendToS3 info_log.txt - ", err)
		}
	} else {
		//fmt.Printf("O arquivo %s não existe na pasta atual.\n", "log.error.txt")
		body, err := awsC.GetInfos("Adm_get_pre_signed?folder=erp_inventory&file_name=info_log")
		if err != nil {
			log.Error(fmt.Sprintf("Log send: GetInfos - %s ", err))
			fmt.Println("BusConfig GetInfos info_log - ", err)
		}

		if err := awsC.SendToS3("info_log.txt", body[strconv.Itoa(1)].(string)); err != nil {
			log.Error(fmt.Sprintf("Log send: SendToS3 - %s ", err))
			fmt.Println("BusConfig SendToS3 info_log.txt - ", err)
		}
	}
	/**/
	log.Info("////////////////////////////////////////////////////////////")
	log.Info(fmt.Sprintf("Proximo envio: %s ", dbInfos["EXECUTION_TIME"].(string)))
	log.Info("////////////////////////////////////////////////////////////")

	return dbInfos["EXECUTION_TIME"].(string)

}
