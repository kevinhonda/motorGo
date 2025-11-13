package diagnosis

import (
	"fmt"
	"motorv2/pkg/db"
	"motorv2/pkg/ws"
	awsC "motorv2/src/awsControllers"

	//"strconv"
	"strings"
	//"bufio"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/kafy11/gosocket/log"
)

func MotorDiagnosis() {
	var macID string

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	ws.SetBaseUrl(os.Getenv("WS_BASE_URL"))
	ws.SetAuth(os.Getenv("WS_AUTH_USER"), os.Getenv("WS_AUTH_PASS"))

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
			macID = macAddr.String()
		}
	}

	//Testando conexao WS
	motorConfigs, err := awsC.GetInfos("Adm_get_wes_config?id=" + macID)
	if err != nil {
		log.Error(fmt.Sprint("GetInfos Error:", err))
		log.Error(fmt.Sprint("Erro de conexão:", err))
		fmt.Println("Erro de conexão:", err)
		fmt.Println("Pressione Enter para sair...")
		fmt.Scanln()
		os.Exit(0)
	}

	//Atribuindo valor para as variaveis de banco
	dbInfos := motorConfigs["env"].(map[string]interface{})

	//Verificando se a chave ["query"] existe no no retorno do serviço Adm_get_wes_config
	//Chave ["query"] é onde vem tanto as querys de envio quanto as de retorno
	if dataQuerys, exists := motorConfigs["query"].(map[string]interface{}); exists {

		//Validando layouts de envio
		if layoutStr, exists := dbInfos["LAYOUTS"].(string); exists {
			layoutStr = strings.ReplaceAll(strings.ToUpper(layoutStr), " ", "")
			log.Info("Layouts de envio: " + layoutStr)
			layouts := strings.Split(layoutStr, ",")
			fmt.Println("Layouts de envio:")
			//Checkpoint
			if newQuery, exists := dataQuerys["NEW"].(map[string]interface{}); exists {
				//Checkpoint
				for i := 0; i < len(layouts); i++ {
					query, ok := newQuery[layouts[i]].(string)

					if ok && query != "" {
						fmt.Println("O layout " + layouts[i] + " possui query")
						log.Info("O layout " + layouts[i] + " possui query")
					} else {
						fmt.Println("O layout " + layouts[i] + " não possui query")
						log.Info("O layout " + layouts[i] + " não possui query")
						log.Error("O layout " + layouts[i] + " não possui query")
					}

				}
			} else {
				fmt.Println("Essa company não possui querys de envio")
				log.Error("Essa company não possui querys de envio")
			}

		} else {
			fmt.Println("Essa company não possui layout de envio")
			log.Error("Essa company não possui layout de envio")
		}

		//Validando layouts de retorno
		if layoutStr, exists := dbInfos["RETURN"].(string); exists {
			layoutStr = strings.ReplaceAll(strings.ToUpper(layoutStr), " ", "")
			log.Info("Layouts de retorno: " + layoutStr)
			layouts := strings.Split(layoutStr, ",")
			fmt.Println("Layouts de retorno:")
			//Checkpoint
			if newQuery, exists := dataQuerys["RETURN"].(map[string]interface{}); exists {
				//Checkpoint
				for i := 0; i < len(layouts); i++ {
					fmt.Println(layouts[i])
					log.Info(layouts[i])
					if query, exists := newQuery[layouts[i]].(map[string]interface{}); exists {
						fmt.Println("O layout " + layouts[i] + " possui querys")
						log.Info("O layout " + layouts[i] + " possui querys")

						queryReturn, ok := query["HEADER"].(string)
						if ok && queryReturn != "" {
							fmt.Println("O layout " + layouts[i] + " possui query de HEADER")
							log.Info("O layout " + layouts[i] + " possui query de HEADER")
						} else {
							fmt.Println("O layout " + layouts[i] + " não possui query de HEADER")
							log.Info("O layout " + layouts[i] + " não possui query de HEADER")
							log.Error("O layout " + layouts[i] + " não possui query de HEADER")
						}

						queryReturn, ok = query["ITEM"].(string)
						if ok && queryReturn != "" {
							fmt.Println("O layout " + layouts[i] + " possui query de ITEM")
							log.Info("O layout " + layouts[i] + " possui query de ITEM")
						} else {
							fmt.Println("O layout " + layouts[i] + " não possui query de ITEM")
							log.Info("O layout " + layouts[i] + " não possui query de ITEM")
							log.Error("O layout " + layouts[i] + " não possui query de ITEM")
						}

						queryReturn, ok = query["PARC"].(string)
						if ok && queryReturn != "" {
							fmt.Println("O layout " + layouts[i] + " possui query de PARC")
							log.Info("O layout " + layouts[i] + " possui query de PARC")
						} else {
							fmt.Println("O layout " + layouts[i] + " não possui query de PARC")
							log.Info("O layout " + layouts[i] + " não possui query de PARC")
							log.Error("O layout " + layouts[i] + " não possui query de PARC")
						}

					} else {
						fmt.Println("O layout " + layouts[i] + " não possui query")
						log.Info("O layout " + layouts[i] + " não possui query")
						log.Error("O layout " + layouts[i] + " não possui query")
					}

				}
			} else {
				fmt.Println("Essa company não possui querys de envio")
				log.Error("Essa company não possui querys de envio")
			}

		} else {
			fmt.Println("Essa company não possui layout de envio")
			log.Error("Essa company não possui layout de envio")
		}
	} else {
		fmt.Println("Essa company não possui querys")
		log.Error("Essa company não possui querys")
	}

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

	fmt.Println("Banco configurado foi: " + dbInfos["DB"].(string))
	log.Info("Banco configurado foi: " + dbInfos["DB"].(string))

	//fmt.Println(motorConfigs["status"].(string))

	mess, dbErr := db.ConnCheck()
	if dbErr != nil {
		log.Error("Erro de conexão com o banco - "+mess, dbErr)

		/*body, err := awsC.GetInfos("Adm_get_pre_signed?folder=erp_inventory&file_name=error_log")
		if err != nil {
			log.Info("Vai passar aquis - awsC")
			fmt.Println(err)
			log.Error(err)
			fmt.Println("Pressione Enter para sair...")
			fmt.Scanln()
			os.Exit(0)
			//return
		}

		if err := awsC.SendToS3("error_log.txt", body[strconv.Itoa(1)].(string)); err != nil {
			log.Error(err)
			fmt.Println(err)
			fmt.Println("Pressione Enter para sair...")
			fmt.Scanln()
			os.Exit(0)
			//return errors.New(fmt.Sprint("Erro ao converter JSON para CSV:", err))
			//return
		}*/

		fmt.Println(mess, dbErr)
		log.Error("Erro de conexão com o banco - "+mess, dbErr)

		fmt.Println("Pressione Enter para sair...")
		fmt.Scanln()
		os.Exit(0)
		//panic(err.Error())
		//return
	}

	log.Info(mess)
	fmt.Println(mess)

	/*
		//Pergunta com resposta aberta
		scanner := bufio.NewScanner(os.Stdin)
		//fmt.Printf("Deseja enviar os layouts manualmente?")
		fmt.Println("Deseja executar qual ação?")
		fmt.Println("1 - Envio")
		fmt.Println("2 - Retorno")


		scanner.Scan()
		//varValue := scanner.Text()
	*/
	fmt.Println("Pressione Enter para sair...")
	fmt.Scanln()
	os.Exit(0)
}
