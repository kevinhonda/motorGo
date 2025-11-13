package handReturn

import (
	"fmt"
	"motorv2/pkg/db"
	"motorv2/src/returnEng"
	"motorv2/src/returnFuncs"

	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"

	"bufio"
	"os"
	"strconv"
	"strings"
)

var dbs *sqlx.DB
var err error

func HandReturn() {
	log.Info("Retorno manual - Inicio")
	fmt.Println("Retorno manual - Inicio")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Deseja executar qual ação?")
	fmt.Println("1 - Retorno")
	fmt.Println("2 - Buss mini")
	fmt.Println("3 - Retorno e Buss mini")
	scanner.Scan()
	varValue := scanner.Text()

	switch varValue {
	case "1":
		for w := 0; w < 1; {
			fmt.Println("Deseja executar qual ação?")
			fmt.Println("1 - Retorno por layout")
			fmt.Println("2 - Retorno de todos os layouts")
			scanner.Scan()
			returnType := scanner.Text()
			if returnType == "1" || returnType == "2" {
				runReturn(returnType)
				w = 1
			} else {
				fmt.Println("Resposta inválida")
				w = 0
			}
		}
	case "2":
		for w := 0; w < 1; {
			fmt.Println("Deseja executar qual ação?")
			fmt.Println("1 - Envio mini por layout")
			fmt.Println("2 - Envio mini de todos os layouts")
			scanner.Scan()
			miniType := scanner.Text()
			if miniType == "1" || miniType == "2" {
				runBussMini(miniType)
				w = 1
			} else {
				fmt.Println("Resposta inválida")
				w = 0
			}
		}

	case "3":
		runBussMini("1")
		runReturn("1")
	default:
		fmt.Println("Resposta inválida")
		fmt.Println("Pressione Enter para sair...")
		fmt.Scanln()
		os.Exit(0)
	}

	fmt.Println("Pressione Enter para sair...")
	fmt.Scanln()
	os.Exit(0)

	fmt.Println("Retorno manual - Final")
	log.Info("Retorno manual - Final")
}

func runBussMini(miniType string) {
	returnFuncsInstance := returnFuncs.GetTest{}
	bodyResp, _ := returnFuncsInstance.GetReturn("Adm_get_wes_config")
	dbInfos := bodyResp["env"].(map[string]interface{})
	dataQuerys := bodyResp["query"].(map[string]interface{})

	//dataQuery := bodyResp["query"].(map[string]interface{})
	log.Info("Querado")
	log.Info("Tem query de LAYOUTS_MINI")
	if returnLayouts, exists := bodyResp["env"].(map[string]interface{})["LAYOUTS_MINI"].(string); exists {
		if returnLayouts != "" {
			if miniType == "1" {
				returnLayouts = perLayout(returnLayouts, "2")
			}

			layoutsMini := strings.ReplaceAll(strings.ToUpper(returnLayouts), " ", "")
			layoutsMiniList := strings.Split(layoutsMini, ",")
			log.Info("Banco configurado foi: " + dbInfos["DB"].(string))

			db.SetConfig(&db.DBConnectionParams{
				DB:     dbInfos["DB"].(string),
				DBUser: dbInfos["DB_USER"].(string),
				DBPass: dbInfos["DB_PASSWORD"].(string),
				DBHost: dbInfos["DB_HOST"].(string),
				DBPort: dbInfos["DB_PORT"].(string),
				DBSid:  dbInfos["DB_SID"].(string),
				DBName: dbInfos["DB_NAME"].(string),
			})

			miniQuerys := make(map[string]string)
			if sendQuery, exists := dataQuerys["NEW"].(map[string]interface{}); exists {
				for _, layoutMini := range layoutsMiniList {
					if sendQueryr, exists := sendQuery[layoutMini].(string); exists && sendQueryr != "" {
						log.Info("Query encontrada para layout:", layoutMini)
						miniQuerys[layoutMini] = sendQueryr
					} else {
						log.Error("Layout sem query:", layoutMini)
					}
				}
			} else {
				log.Error("Sem querys para a company")
				fmt.Println("Sem querys para a company")
			}
			if len(miniQuerys) > 0 {
				//log.Info("runBussMini - Querys Layout Mini")
				log.Info(miniQuerys)
				//---------------------------------------------------------------------------------
				if dbs, err = db.OpenDB(); err != nil {
					fmt.Println("Error initializing database connection:", err)
					log.Error("Error initializing database connection:", err)
				}

				defer func() {
					err := dbs.Close()
					if err != nil {
						log.Error("Erro ao fechar a conexão com o banco de dados:", err)
					}
					fmt.Println("A conexão com o banco de dados foi fechada.")
					log.Info("ReturnAction - A conexão com o banco de dados foi fechada.")
				}()
				//---------------------------------------------------------------------------------
				returnEng.RunMinis(miniQuerys)
			} else {
				fmt.Println("Sem querys para layouts mini")
				log.Info("Sem querys para layouts mini")
			}

		} else {
			fmt.Println("A company não tem layouts mini")
			log.Info("A company não tem layouts mini")
		}
	}
}

func runReturn(returnType string) {
	returnFuncsInstance := returnFuncs.GetTest{}
	bodyResp, _ := returnFuncsInstance.GetReturn("Adm_get_wes_config")
	dbInfos := bodyResp["env"].(map[string]interface{})
	dataQuery := bodyResp["query"].(map[string]interface{})

	if returnLayouts, exists := bodyResp["env"].(map[string]interface{})["RETURN"].(string); exists {
		returnLayouts := strings.ReplaceAll(strings.ToUpper(returnLayouts), " ", "")
		if returnLayouts != "" {
			if returnType == "1" {
				//Pergunta com resposta aberta
				//layouts := strings.Split(returnLayouts, ",")

				returnLayouts = perLayout(returnLayouts, "1")
			}

			log.Info("Banco configurado foi: " + dbInfos["DB"].(string))

			db.SetConfig(&db.DBConnectionParams{
				DB:     dbInfos["DB"].(string),
				DBUser: dbInfos["DB_USER"].(string),
				DBPass: dbInfos["DB_PASSWORD"].(string),
				DBHost: dbInfos["DB_HOST"].(string),
				DBPort: dbInfos["DB_PORT"].(string),
				DBSid:  dbInfos["DB_SID"].(string),
				DBName: dbInfos["DB_NAME"].(string),
			})

			//---------------------------------------------------------------------------------
			if dbs, err = db.OpenDB(); err != nil {
				fmt.Println("Error initializing database connection:", err)
				log.Error("Error initializing database connection:", err)
			}

			defer func() {
				err := dbs.Close()
				if err != nil {
					log.Error("Erro ao fechar a conexão com o banco de dados:", err)
				}
				fmt.Println("A conexão com o banco de dados foi fechada.")
				log.Info("ReturnAction - A conexão com o banco de dados foi fechada.")
			}()
			//---------------------------------------------------------------------------------
			if dataQuery, exists := dataQuery["RETURN"].(map[string]interface{}); exists {
				log.Info("returnQuery:")
				log.Info("ReturnAction - sendQuery")
				log.Info("returnQuerys:")

				returnEng.RunReturns(returnLayouts, dataQuery)
			} else {
				log.Info("dataQuery RETURN - Dequeryado")
				//returnEng.RunReturns(returnLayouts, nil)
			}

			log.Info("Executou o run returns")
		} else {
			log.Info("Essa company não possui layout de retorno")
			fmt.Println("Essa company não possui layout de retorno")
		}
	} else {
		log.Error("Não possui query para a company ou para o layout.")
	}
}

func perLayout(layoutNames, typeReturn string) string {
	scanner := bufio.NewScanner(os.Stdin)

	ReturnName := ""
	switch typeReturn {
	case "1":
		ReturnName = "retorno"
	case "2":
		ReturnName = "buss mini"
	default:
		log.Error("perLayout: Tipo invalido")
	}

	log.Info(fmt.Sprintf("Layouts para %s: %s", ReturnName, layoutNames))
	fmt.Printf("Layouts para %s: %s \n", ReturnName, layoutNames)
	layoutRtnStr := strings.ReplaceAll(strings.ToUpper(layoutNames), " ", "")
	log.Info("Layouts: " + layoutRtnStr)
	layouts := strings.Split(layoutRtnStr, ",")

	if len(layouts) > 1 {
		fmt.Printf("Tem algo ais \n")
	} else {
		fmt.Printf("Tem nada não bicho \n")
	}

	var layoutsName string
	for w := 0; w < 1; {
		for index, layout := range layouts {
			fmt.Printf("%d - %s\n", index, layout)
		}
		fmt.Printf("Quais layouts deseja rodar %s\n", ReturnName)
		scanner.Scan()
		layoutsValues := scanner.Text()
		layoutsValues = strings.ReplaceAll(strings.ToUpper(layoutsValues), " ", "")
		layoutsVavalue := strings.Split(layoutsValues, ",")
		layoutsName = ""
		w = len(layoutsVavalue)
		for _, layoutIndex := range layoutsVavalue {
			// Convert string to integer
			index, err := strconv.Atoi(layoutIndex)
			if err != nil {
				fmt.Println("///////////////////////////////////////")
				fmt.Printf("Erro: '%s' não é um número válido\n", layoutIndex)
				fmt.Println("///////////////////////////////////////")
				w = w - 1
				continue
			}

			// Check if index is within bounds
			if index < 0 || index >= len(layouts) {
				fmt.Println("///////////////////////////////////////")
				fmt.Printf("Erro: índice %d fora dos limites (0-%d)\n", index, len(layouts)-1)
				fmt.Println("///////////////////////////////////////")
				w = w - 1
				continue
			}

			layoutsName += layouts[index]
			layoutsName += ", "
		}

		if len(layoutsName) > 0 {
			layoutsName = layoutsName[:len(layoutsName)-2]
		}

		//fmt.Println("Layouts para retorno: " + layoutsName)
		fmt.Printf("Deseja executar o %s desses layous %s\n", ReturnName, layoutsName)
		//fmt.Println("Deseja executar o retorno desses layous " + layoutsName)
		if w > 0 {
			for iw := 0; iw < 1; {
				scanner.Scan()
				varValue := strings.ToLower(scanner.Text())

				switch varValue {
				case "sim", "s":
					w = 1
					iw = 2
				case "não", "n", "nao":
					w = -1
					iw = 2
				default:
					//fmt.Println("Resposta inválida. Por favor, responda com 'sim', 's', 'não' ou 'n'.")
					fmt.Println("////////////////////////////////////")
					fmt.Println("//////////Resposta inválida/////////")
					fmt.Println("////////////////////////////////////")
				}
			}
		}
	}
	return layoutsName

}
