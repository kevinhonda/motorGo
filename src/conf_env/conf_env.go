package conf_env

import (
	"bufio"

	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"motorv2/pkg/db"
	"motorv2/pkg/ws"
	awsC "motorv2/src/awsControllers"

	"github.com/kafy11/gosocket/log"
)

var baseUrl string

func SetUrl(params string) {
	switch params {
	case "QA":
		baseUrl = "https://api-qa.gtplan.net"
	case "PRD":
		baseUrl = "https://api.gtplan.net"
	default:
		baseUrl = "https://api-qa.gtplan.net"
	}
	log.Info("Base URL set to: " + baseUrl)
}

/*
	func loadEnvFile(filePath string) (map[string]string, error) {
		err := godotenv.Load(filePath)
		if err != nil {
			return nil, err
		}

		envMap := make(map[string]string)

		for _, envVar := range os.Environ() {
			keyValue := strings.SplitN(envVar, "=", 2)
			if len(keyValue) == 2 {
				envMap[keyValue[0]] = keyValue[1]
			}
		}

		return envMap, nil
	}
*/
func getdb(dbname string) map[int]string {
	outerMap := make(map[string]map[int]string)

	oracle := map[int]string{}
	oracle[1] = "DB_USER"
	oracle[2] = "DB_PASSWORD"
	oracle[3] = "DB_HOST"
	oracle[4] = "DB_PORT"
	oracle[5] = "DB_SID"
	//oracle[6] = "DB_NAME"
	/*
		postgre := map[int]string{}
		postgre[1] = "DB_USER"
		postgre[2] = "DB_PASSWORD"
		postgre[3] = "DB_HOST"
		postgre[4] = "DB_PORT"
		postgre[5] = "DB_NAME"
	*/
	postgre := map[int]string{}
	postgre[1] = "DB_USER"
	postgre[2] = "DB_PASSWORD"
	postgre[3] = "DB_NAME"

	sqlserver := map[int]string{}
	sqlserver[1] = "DB_USER"
	sqlserver[2] = "DB_PASSWORD"
	sqlserver[3] = "DB_HOST"
	sqlserver[4] = "DB_NAME"

	outerMap["Oracle"] = oracle
	outerMap["Postgre"] = postgre
	outerMap["SQLserver"] = sqlserver

	//newdbtype := map[string]string{}

	newdbtype := outerMap[dbname]
	return newdbtype
}

func mergeMaps(map1, map2 map[string]string) map[string]string {
	result := make(map[string]string)

	// Adiciona todos os elementos do primeiro mapa
	for key, value := range map1 {
		result[key] = value
	}

	// Adiciona todos os elementos do segundo mapa
	for key, value := range map2 {
		result[key] = value
	}

	return result
}

func validateParams() int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		//fmt.Println("Os parâmetros estão corretos? (sim/s/não/n)")
		fmt.Println("Os parâmetros estão certo?")
		scanner.Scan()
		varValue := strings.ToLower(scanner.Text())

		switch varValue {
		case "sim", "s":
			return 1
		case "não", "n":
			return -1
		default:
			//fmt.Println("Resposta inválida. Por favor, responda com 'sim', 's', 'não' ou 'n'.")
			fmt.Println("////////////////////////////////////")
			fmt.Println("//////////Resposta inválida/////////")
			fmt.Println("////////////////////////////////////")
		}
	}
}

func getParams(pergunta string, optiones int, values map[int]string) (map[string]string, string) {
	varValue := ""
	scanner := bufio.NewScanner(os.Stdin)
	switch optiones {
	case 1:
		//Values como opções
		//Opções
		opt := len(values)
		//Exibindo pergunta
		fmt.Println(pergunta)
		log.Info(pergunta)
		//Loop para não deixar usuario escolher algo fora das opções
		for w := 0; w < 1; {
			/*for key, value := range values {
				fmt.Printf("%s - %s\n", key, value)
			}*/
			//Exibindo as opções
			for i := 1; i < len(values)+1; i++ {
				fmt.Printf("%d - %s\n", i, values[i])
			}
			//Pegando opção escolhida
			fmt.Scanln(&varValue)
			log.Info(varValue)
			//Verificando se usuario botou exit
			strLower := strings.ToLower(varValue)
			if strLower == "exit" {
				os.Exit(0)
				return nil, "exit"
			}
			//Validando a opções escolhida
			validOpt, err := strconv.Atoi(varValue)
			if err != nil {
				fmt.Println("Erro ao converter string para int:", err)
				w = 0
			}
			if opt < validOpt || validOpt < 1 {
				fmt.Println("Opção inválida")
				w = 0
			} else {
				w = 1
			}
		}
		varValue, err := strconv.Atoi(varValue)
		if err != nil {
			fmt.Println("Erro ao converter string para int:", err)
		}

		log.Info(pergunta + " " + values[varValue])
		return nil, values[varValue]

	case 2:
		//Values sendo preenchido
		//Map a ser preenchido
		mapFill := make(map[string]string)
		//outerMap := make(map[string]map[int]string)
		for i := 1; i < len(values)+1; i++ {
			//Pergunta que vai ser exibida
			question := fmt.Sprintf("%s %s: ", pergunta, values[i])
			fmt.Printf(question)
			log.Info(question)
			//Lendo retorno do usuario
			scanner.Scan()
			varValue := scanner.Text()
			log.Info(varValue)
			//Verificando se usuario botou exit
			strLower := strings.ToLower(varValue)
			if strLower == "exit" {
				os.Exit(0)
				return nil, "exit"
			}
			//Preenchendo o map
			mapFill[values[i]] = varValue
		}
		return mapFill, "nice"

	case 3:
		//Pergunta com resposta aberta
		//Exibindo Pergunta
		fmt.Printf("%s ", pergunta)
		log.Info(pergunta)
		//Pegando resposta do usuario
		scanner.Scan()
		varValue := scanner.Text()
		log.Info(pergunta)
		//Verificando se usuario botou exit
		strLower := strings.ToLower(varValue)
		if strLower == "exit" {
			os.Exit(0)
			return nil, "exit"
		}
		log.Info(pergunta + " " + varValue)
		return nil, varValue
	}
	return nil, "exit"
}

func timeValidate(executeTimes []string) (string, int) {
	for i := 0; i < len(executeTimes); i++ {
		executeTimes[i] = strings.ReplaceAll(executeTimes[i], " ", "")
		if !strings.Contains(executeTimes[i], ":") {
			log.Error("Falta : no Horario")
			return fmt.Sprint("Falta : no Horario"), -1
		}
		timeParts := strings.Split(executeTimes[i], ":")
		if len(timeParts) != 2 {
			log.Error("A variável de ambiente EXECUTION_TIME não está no formato correto (HH:MM).")
			return fmt.Sprint("A variável de ambiente EXECUTION_TIME não está no formato correto (HH:MM)."), -1
		}
		executionHour, err := strconv.Atoi(timeParts[0])
		if err != nil {
			log.Error(fmt.Sprint("Erro ao analisar a hora de execução: ", err))
			return fmt.Sprint("Erro ao analisar a hora de execução: ", err), -1
		}
		executionMinute, err := strconv.Atoi(timeParts[1])
		if err != nil {
			log.Error(fmt.Sprint("Erro ao analisar a minutos de execução: ", err))
			return fmt.Sprint("Erro ao analisar a minutos de execução: ", err), -1
		}

		if executionHour > 23 || executionHour < 0 || executionMinute > 59 || executionMinute < 0 {
			return fmt.Sprint("Horas ou minutos incorreto: ", executeTimes[i]), -1
		}

		regex := regexp.MustCompile(`^\d{2}:\d{2}$`)
		if regex.MatchString(executeTimes[i]) {
			log.Info("Formato de horário correto:", executeTimes[i])
		} else {
			return fmt.Sprint("Formato de horário incorreto:", executeTimes[i]), -1
		}
	}
	return "Horarios valido", 5
}

func SetConf() {
	envVariables := map[string]string{
		"WS_BASE_URL":  "",
		"WS_AUTH_USER": "",
		"WS_AUTH_PASS": "",
	}

	variables := map[string]string{
		"ID":                   "",
		"DB":                   "",
		"ERP":                  "",
		"ID_ERP":               "",
		"COMPANY_NAME":         "",
		"EXECUTION_TIME":       "",
		"DRIVER_LETTER":        "C",
		"CONTROLLERS_DIR_PATH": "./src/logAll",
		"SERVICE_NAME":         ".GTPMotorV2",
		"WEBSOCKET_SERVER_URL": "",
		//"RETURN": "Mrp_po, Sug_id, Drp_id",
		"RETURN":       "",
		"LAYOUTS_MINI": "",
		"LAYOUTS":      "",
		//"LAYOUTS": "ERP_INVENTORY, ERP_RESERVE, ERP_SUPPLIER,ERP_SKU, ERP_SKU_ALL, ERP_PURCHASE_REQ",
	}

	//"RETURN":  "Mrp_po, Sug_id, Drp_id, Surgen",
	//"LAYOUTS":              "erp_inventory,ERP_PURCHASE_TESTER,ERP_TESTER_ORDER,ERP_TESTERU",

	dbParams := map[string]string{
		"DB_USER":     "",
		"DB_PASSWORD": "",
		"DB_HOST":     "",
		"DB_PORT":     "",
		"DB_SID":      "",
		"DB_NAME":     "",
	}

	erps := map[int]string{
		1: "MV",
		2: "Tasy",
		3: "Outro",
		4: "TESTER",
	}

	idErps := map[string]string{
		"MV":     "-1",
		"Tasy":   "-2",
		"Outro":  "-3",
		"TESTER": "-10",
	}

	dbs := map[int]string{
		1: "Oracle",
		2: "SQLserver",
		3: "Postgre",
	}
	/*
		ambient := map[int]string{
			1: "QA",
			2: "PROD",
			3: "D1",
			4: "IND",
		}

		links := map[string]string{
			"QA":   "https://api-qa.gtplan.net",
			"PROD": "gtplanosprod.cooler",
			"D1":   "gtplanosd1.cooler",
			"IND":  "gtplanosind.cooler",
		}

	*/
	//Verificando se o arquivo .env existe
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("O arquivo .env não existe no diretório atual. Criando o arquivo...")

		// Criar o arquivo .env
		file, err := os.Create(".env")
		if err != nil {
			fmt.Println("Erro ao criar o arquivo .env:", err)
			return
		}
		defer file.Close()
		fmt.Println("Arquivo .env criado com sucesso!")
	}

	dir, err := os.Getwd()
	// Verifica se ocorreu algum erro
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho:", err)
		os.Exit(1)
	}

	fmt.Println("Diretório de trabalho atual:", dir)
	fmt.Println("Configuração do arquivo .env:")
	fmt.Println("Digite exit em qualquer momento para encerrar.")

	envVariables["WS_BASE_URL"] = baseUrl

	fmt.Println("Link:", envVariables["WS_BASE_URL"])
	for {
		_, varValue := getParams("AUTH USER:", 3, nil)
		envVariables["WS_AUTH_USER"] = varValue

		_, varValue = getParams("AUTH PASS:", 3, nil)
		envVariables["WS_AUTH_PASS"] = varValue

		validationResult := validateParams()

		if validationResult == 1 {
			break
		}
	}
	//////////////////////////////////////////////////////////////////////////
	//Configuração BANCO
	//////////////////////////////////////////////////////////////////////////

	ws.SetBaseUrl(envVariables["WS_BASE_URL"])
	ws.SetAuth(envVariables["WS_AUTH_USER"], envVariables["WS_AUTH_PASS"])
	for {
	loopConfig:
		for {
			//Pegando valores do banco
			_, dbname := getParams("Escolha uma opção de banco:", 1, dbs)
			variables["DB"] = dbname

			dbFields := getdb(dbname)

			dbValues, _ := getParams("Digite o valor do campo -", 2, dbFields)
			log.Info("O banco escolhido foi: " + variables["DB"])

			//log.Info("dbValues")
			//log.Info(dbValues)

			for key := range dbValues {
				dbParams[key] = dbValues[key]
			}

			//log.Info("dbParams")
			//log.Info(dbParams)

			validationResult := validateParams()
			if validationResult == 1 {
				break loopConfig
			}
		}

		//Setando valores banco
		db.SetConfig(&db.DBConnectionParams{
			DB:     variables["DB"],
			DBUser: dbParams["DB_USER"],
			DBPass: dbParams["DB_PASSWORD"],
			DBHost: dbParams["DB_HOST"],
			DBPort: dbParams["DB_PORT"],
			DBSid:  dbParams["DB_SID"],
			DBName: dbParams["DB_NAME"],
		})

		validationResult := 1
		//Teste de banco
		mess, dbErr := db.ConnCheck()

		//Se erro banco
		if dbErr != nil {
			log.Error("Configuração de banco - "+mess, dbErr)
			fmt.Println(mess, dbErr)
			log.Error("Configuração de banco - "+mess, dbErr)
		loopReconfig:
			for {
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Println("A conexão deu errado.")
				fmt.Println("Gostaria de tentar novamente?")
				scanner.Scan()
				varValue := strings.ToLower(scanner.Text())

				switch varValue {
				case "sim", "s":
					validationResult = -1
					break loopReconfig
				case "não", "n":
					fmt.Println("Pressione Enter para sair...")
					fmt.Scanln()
					os.Exit(0)
				default:
					//fmt.Println("Resposta inválida. Por favor, responda com 'sim', 's', 'não' ou 'n'.")
					fmt.Println("////////////////////////////////////")
					fmt.Println("//////////Resposta inválida/////////")
					fmt.Println("////////////////////////////////////")
				}
			}
		}

		//Log
		log.Info(mess)
		fmt.Println(mess)

		//Validação Loop
		if validationResult == 1 {
			break
		}

	}
	//////////////////////////////////////////////////////////////////////////
	//Configuração ERP/Horario de execução/Company Name
	//////////////////////////////////////////////////////////////////////////
	for {
		//Trecho conf de ERP
		_, varValue := getParams("Escolha uma opção de ERP:", 1, erps)
		variables["ERP"] = varValue
		variables["ID_ERP"] = idErps[varValue]
		if variables["ERP"] == "Outro" {
			_, varValue = getParams("Qual o nome do ERP", 3, erps)
			variables["ERP"] = varValue
		}
		//Trecho conf de Company e horario
		_, varValue = getParams("COMPANY NAME:", 3, nil)
		variables["COMPANY_NAME"] = varValue
		for w := 0; w < 1; w++ {
			var msg string
			_, varValue = getParams("Horario de envio no formato (HH:MM):", 3, nil)
			variables["EXECUTION_TIME"] = varValue
			executeTimes := strings.Split(variables["EXECUTION_TIME"], ",")
			msg, w = timeValidate(executeTimes)
			fmt.Println(msg)
		}

		validationResult := validateParams()
		if validationResult == 1 {
			break
		}
	}

	log.Info("Time:" + variables["EXECUTION_TIME"])

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Erro ao buscar interfaces de rede:", err)
		return
	}

	// Loop pelas interfaces para encontrar o endereço MAC
	for _, iface := range interfaces {
		// Ignora interfaces loopback e interfaces sem endereço MAC
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			macAddr := iface.HardwareAddr
			fmt.Println("Endereço MAC:", macAddr.String())
			log.Info("Endereço MAC:" + macAddr.String())
			variables["ID"] = macAddr.String()
		}
	}

	//Criando array de variaveis q vão ser enviadas para o dynamo
	dynamoVariables := mergeMaps(variables, dbParams)
	dynamoVariables = mergeMaps(dynamoVariables, envVariables)

	log.Info(envVariables["WS_BASE_URL"])
	log.Info("Variaveis")
	log.Info(dynamoVariables)
	//Enviando para o dynamo
	errSP := awsC.SendParams(dynamoVariables, envVariables["WS_BASE_URL"], envVariables["WS_AUTH_USER"], envVariables["WS_AUTH_PASS"])
	if errSP != "" {
		fmt.Println("Erro ao mandar env:", errSP)
		//log.Error("Erro ao mandar env:", err)
		log.Error(fmt.Sprintf("%s - Erro ao mandar env: %s", envVariables["WS_BASE_URL"], errSP))
		fmt.Println("Pressione Enter para sair...")
		fmt.Scanln()
		os.Exit(0)
		//return
	}

	// Escreve as variáveis no arquivo .env
	envFile, err := os.Create(".env")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo .env:", err)
		return
	}
	defer envFile.Close()

	for key, value := range envVariables {
		_, err := fmt.Fprintf(envFile, "%s=%s\n", key, value)
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo .env:", err)
			return
		}
	}
	//imprimirMapa(envVariables)
	fmt.Println("Arquivo .env configurado com sucesso!")
	/*
		body, err := awsC.GetInfos("Adm_get_pre_signed?folder=erp_inventory&file_name=info_log")
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			return
		}

		if err := awsC.SendToS3("info_log.txt", body[strconv.Itoa(1)].(string)); err != nil {
			log.Error(err)
			fmt.Println(err)
			return
		}
	*/

}
