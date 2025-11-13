package returnFuncs

import (
	"encoding/json"
	"fmt"

	"sort"
	"time"

	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"runtime"

	"github.com/jmoiron/sqlx"

	"motorv2/pkg/db"
	"motorv2/pkg/ws"

	"github.com/kafy11/gosocket/log"
)

type GetTest struct{}

func GetReturns(layoutReturn, controller string) string {

	var returnLayoutName string
	switch layoutReturn {
	case "Mrp_po":
		returnLayoutName = "Ordem"
	case "Sug_id":
		returnLayoutName = "Solicitacao"
	case "Drp_id":
		returnLayoutName = "Transferencia"
	default:
		returnLayoutName = "default"
	}

	//log.Info("OC - Antes de rodar")
	returnFuncsInstance := GetTest{}
	body, err := returnFuncsInstance.GetReturn(controller)
	if err != nil {
		//log.Error(fmt.Sprint("GetReturn Error:", err))
		//fmt.Println(fmt.Sprint("GetReturn Error:", err))
		return fmt.Sprint("GetReturn Error:", err)
	}

	fmt.Println("O nome do layout é - ", layoutReturn)
	//jsonsBody := body["RESULT"].(map[string]interface{})["DATA"].(map[string]interface{})[returnLayoutName].([]interface{})
	var jsonsBody ([]interface{})
	if result, ok := body["RESULT"].(map[string]interface{}); ok {
		// Verificar se a chave "DATA" existe dentro de "RESULT"

		if returnLayoutName != "default" {

			if data, exists := result["DATA"].(map[string]interface{}); exists {
				fmt.Println("A chave 'DATA' existe")
				log.Info(controller, " - A chave 'DATA' existe")

				jsonsBody = data[returnLayoutName].([]interface{})

				folderName := "jsonFiles"
				if _, err := os.Stat(folderName); os.IsNotExist(err) {
					fmt.Println("Directory does not exist.")
					log.Error("Directory does not exist.")

					err := os.Mkdir(folderName, 0755)
					if err != nil {
						//fmt.Println("Erro ao criar a pasta:", err)
						//log.Error("Erro ao criar a pasta:", err)
						return fmt.Sprint("Erro ao criar a pasta:", err)
					}
				} else {
					fmt.Println("Directory exists.")
				}
				//Criando JSONs na pasta
				CreateJson(layoutReturn, folderName, jsonsBody)

			} else {
				//fmt.Println("A chave 'DATA' não existe.")
				//log.Error(controller, " - A chave 'DATA' não existe.")
				return fmt.Sprintf("%s - A chave 'DATA' não existe.", controller)
			}

		} else {
			if data, exists := result["DATA"].([]interface{}); exists {
				fmt.Println("A chave 'DATA' existe")
				log.Info(controller, " - A chave 'DATA' existe")

				jsonsBody = data

				folderName := "jsonFiles"
				if _, err := os.Stat(folderName); os.IsNotExist(err) {
					fmt.Println("Directory does not exist.")
					log.Error("Directory does not exist.")

					err := os.Mkdir(folderName, 0755)
					if err != nil {
						//fmt.Println("Erro ao criar a pasta:", err)
						//log.Error("Erro ao criar a pasta:", err)
						return fmt.Sprint("Erro ao criar a pasta:", err)
					}
				} else {
					fmt.Println("Directory exists.")
				}
				//Criando JSONs na pasta
				CreateJson(layoutReturn, folderName, jsonsBody)

			} else {
				//fmt.Println("A chave 'DATA' não existe.")
				//log.Error(controller, " - A chave 'DATA' não existe.")
				return fmt.Sprintf("%s - A chave 'DATA' não existe.", controller)
			}
		}

	} else {
		//fmt.Println(controller, " - A chave 'RESULT' não existe ou não é do tipo esperado.")
		//log.Error(controller, " - A chave 'RESULT' não existe ou não é do tipo esperado.")
		return fmt.Sprintf("%s - A chave 'RESULT' não existe ou não é do tipo esperado.", controller)
	}
	log.Info(layoutReturn, " - Terminou a parada")
	return ""
}

func (ra GetTest) GetReturn(controller string) (map[string]interface{}, error) {
	//Realizando chamada
	body, err := ws.Get(controller)
	if err != nil {
		log.Error("GetReturn - Body:", err)
		return nil, err
	}
	//fmt.Printf("O tipo da variável bodys é: %T\n", body)

	var dataBody map[string]interface{}
	// Faz o parsing do JSON
	err = json.Unmarshal(body, &dataBody)
	if err != nil {
		log.Error("GetReturn - Erro ao fazer o parsing do JSON:", err)
		return nil, fmt.Errorf("Erro ao fazer o parsing do JSON: %w", err)
	}
	//log.Info("dataBody:")
	//log.Info(dataBody)
	return dataBody, nil
}

func CreateJson(jsonName string, folderName string, jfiles []interface{}) {

	// Dados para serem escritos no JSON
	count := len(jfiles)

	for i := 0; i < count; i++ {
		// Abrir o arquivo dentro da pasta
		//file, err := os.Create(folderName + "/data.json")
		file, err := os.Create(folderName + "/" + jsonName + "_" + strconv.Itoa(i) + ".json")
		if err != nil {
			fmt.Println("Erro ao criar o arquivo JSON:", err)
			log.Error("Erro ao criar o arquivo JSON:", err)
			return
		}
		defer file.Close()

		// Criar um codificador JSON
		encoder := json.NewEncoder(file)

		// Codificar os dados para JSON e escrever no arquivo
		err = encoder.Encode(jfiles[i])
		if err != nil {
			fmt.Println("Erro ao codificar dados para JSON:", err)
			log.Error("Erro ao codificar dados para JSON:", err)
			return
		}

		//fmt.Println("Pasta e arquivo JSON foram criados com sucesso.")
	}
}

func JsonToArray(jsonName string) map[string]interface{} {
	file, err := os.Open("jsonFiles/" + jsonName)
	log.Info("JsonToArray - Abrindo arquivo:", "jsonFiles/"+jsonName)
	log.Info(file)
	if err != nil {
		fmt.Println("Error:", err)
		log.Error("JsonToArray - Error:", err)
	}
	defer file.Close()

	// Decode JSON from the file
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error:", err)
		log.Error("JsonToArray - Error:", err)
	}

	return data
}

func GetFilesName(folderName, layout string) []string {
	// Ler o diretório
	files, err := os.ReadDir(folderName)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("GetFiles - ", err)
		log.Error("GetFiles - ", err)

	}

	var fileNames []string
	for _, file := range files {
		// Ignorar diretórios
		if file.IsDir() {
			continue
		}
		// Verificar se o nome do arquivo contém a substring desejada
		if strings.Contains(file.Name(), layout) {
			fileNames = append(fileNames, file.Name())
		}
	}

	fmt.Println("Files in", folderName)
	/*
		for _, name := range fileNames {
			fmt.Println(name)
		}
	*/
	return fileNames
}

func DeleteFolder() {
	//Apagar JSONs da pasta jsonFiles
	folderPath := "jsonFiles"
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Erro ao ler a pasta:", err)
		return
	}

	for _, file := range files {
		filePath := folderPath + "/" + file.Name()
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Erro ao apagar o arquivo:", filePath, err)
		} else {
			fmt.Println("Arquivo apagado:", filePath)
		}
	}
}

//var tx *sqlx.Tx

var dbs *sqlx.DB

type ReturnInserts struct {
	db *sqlx.DB
}

func NewReturnInserts(db *sqlx.DB) *ReturnInserts {
	return &ReturnInserts{db: db}
}

//type ReturnInserts struct{}

func (ra ReturnInserts) HeaderInsert(data map[string]interface{}, filesName []string, varStuffs map[string]string) map[string]interface{} {
	//Contando quantos arquivos temos para o retorno
	count := len(filesName)
	queryPoronoue := data["HEADER"].(string)

	myMap := make(map[string]interface{})
	//myMap["Ordem"] = []interface{}{}
	myMap[varStuffs["myMap"]] = []interface{}{}
	log.Info("Header - Chamando db.OpenTX")
	/*
		tx, err := db.OpenTX()
		if err != nil {
			log.Error("Erro ao fazer OpenTX:", err)
		}

		// Garante rollback se algo der errado
		var txErr error
		defer func() {
			if txErr != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Error("Erro ao fazer rollback:", rbErr)
				}
			}
		}()*/
	//Loop Inicial para inserir o cabeçalho do retorno
	//log.Info(varStuffs["myMap"] + "            " + queryPoronoue)
	for i := 0; i < count; i++ {
		// Inicia a transação
		tx, err := db.OpenTX()
		log.Info("Header - Chamando db.OpenTX2")
		if err != nil {
			log.Error("Erro ao iniciar transação:", err)
		}
		log.Info("Abrindo o arquivo: " + filesName[i])
		// Imprimir algumas estatísticas de uso de memória
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		//log.Info("Uso de memória: %d bytes\n", memStats.Alloc)
		//log.Info("Uso de memória total: %d bytes\n", memStats.TotalAlloc)
		//log.Info("Uso de memória do sistema: %d bytes\n", memStats.Sys)
		//log.Info("Número de alocações: %d\n", memStats.Mallocs)
		//log.Info("Número de liberações: %d\n", memStats.Frees)

		cdStatus := "3"
		erroLog := ""
		fileArray := JsonToArray(filesName[i])

		//PO APENAS
		if varStuffs["myMap"] == "Ordem" {
			items := fileArray["ITENS"].([]interface{})
			itemsCount := len(items)
			totalVal := 0.0
			for i := 0; i < itemsCount; i++ {
				itemes := items[i].(map[string]interface{})

				itemes["cdSolicitacao"] = fileArray["cdSolicitacao"]

				qtMaterial, err := strconv.ParseFloat(itemes["qtMaterial"].(string), 64)
				if err != nil {
					fmt.Println("Erro na conversão:", err)
				}

				vlUnitario, err := strconv.ParseFloat(itemes["vlUnitario"].(string), 64)
				if err != nil {
					fmt.Println("Erro na conversão:", err)
				}

				itemVal := qtMaterial * vlUnitario

				totalVal += itemVal
			}
			fileArray["valTotal"] = totalVal
		}
		log.Info("Header - Chamando db.SqlExec")
		resultado, errMsg := db.SqlExec(tx, queryPoronoue, fileArray)
		fmt.Println("Header - ID gerado:", resultado)
		//log.Info(fmt.Sprintf("Executando o header %s", fileArray["nrDocumentoExterno"]))
		log.Info("Header - ID gerado:", resultado)
		log.Info(fmt.Sprintf("Executando o header %s", fileArray[varStuffs["fileArray"]]))

		if errMsg != "" {
			//erroLog = fmt.Sprintf("Erro no nrDocumentoExterno %s - %s", fileArray["nrDocumentoExterno"], errMsg)
			erroLog = fmt.Sprintf("Erro no %s %s - %s", varStuffs["fileArray"], fileArray[varStuffs["fileArray"]], errMsg)
			fmt.Println(erroLog)
			log.Error(erroLog)
			erroLog = errMsg
			cdStatus = "4"
			//tx.Rollback()
		}

		//fmt.Printf("O tipo da variável fileArray é: %T\n", fileArray)

		fileArray["ID_HEADER_STUFF"] = resultado

		ordem := map[string]interface{}{
			varStuffs["fileArray"]: "",
			"cdStatus":             cdStatus,
			varStuffs["header"]:    "",
			"dsLogError":           erroLog,
			"Itens":                []interface{}{},
		}
		ordem[varStuffs["fileArray"]] = fileArray[varStuffs["fileArray"]].(string)
		ordem[varStuffs["header"]] = strconv.Itoa(resultado)

		///////////////////////////
		////////Loop Itens/////////
		///////////////////////////
		//[0].(map[string]interface{})
		var varItens string
		items := []interface{}{}
		if varStuffs["myMap"] == "Transferencia" {
			items = fileArray["Itens"].([]interface{})
			varItens = "Itens"
		} else {
			items = fileArray["ITENS"].([]interface{})
			varItens = "ITENS"
		}
		//filesNome := filesName[i]
		//log.Info(fmt.Sprintf("%s tem %d itens a serem paranauados", filesNome, itemsCount))

		////////////////////////////
		var campos []string
		for key := range fileArray {
			campos = append(campos, key)
		}

		itemsCount := len(items)
		for i := 0; i < itemsCount; i++ {
			itemes := items[i].(map[string]interface{})

			for c := 0; c < len(campos); c++ {
				if campos[c] != varItens {
					itemes[campos[c]] = fileArray[campos[c]]
				}
			}
		}
		////////////////////////////

		itensError := itemInsert(tx, items, data, ordem, resultado, varStuffs)
		if itensError != "" {
			log.Error(itensError)
			ordem["dsLogError"] = "Erro nos itens"
			ordem["cdStatus"] = "4"
		}
		myMap[varStuffs["myMap"]] = append(myMap[varStuffs["myMap"]].([]interface{}), ordem)

		//////////////////////////////
		//Alterar arquivo JSON
		//////////////////////////////
		updatedJSON, err := json.MarshalIndent(fileArray, "", "    ")
		if err != nil {
			fmt.Println("Erro ao codificar o JSON:", err)
			log.Error("Erro ao codificar o JSON:", err)
			//return
		}

		// Escrever os dados de volta no arquivo JSON
		err = ioutil.WriteFile("jsonFiles/"+filesName[i], updatedJSON, 0644)
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			log.Error("Erro ao escrever no arquivo:", err)
			//return
		}

		fmt.Println(fmt.Sprintf("Arquivo JSON %s atualizado com sucesso.", filesName[i]))

		//log.Info("Terminado essa header ")
		log.Info(fmt.Sprintf("Finalizado header %s", fileArray["nrDocumentoExterno"]))
		if cdStatus == "4" {
			tx.Rollback()
			log.Error("Erro no cabeçalho, rollback")
		} else {
			// Commit the transaction
			log.Info("Header - Commitando transação")
			if err := tx.Commit(); err != nil {
				log.Error("Erro ao commitar:", err)
			}
		}
	}

	return myMap
}

func itemInsert(tx *sqlx.Tx, items []interface{}, data map[string]interface{}, ordem map[string]interface{}, idHeader int, varStuffs map[string]string) string {
	itemsCount := len(items)
	log.Info("CD STATUS DESSA HEADER:")
	log.Info(ordem["cdStatus"])
	log.Info("Total de itens: ", itemsCount)
	queryItem := data["ITEM"].(string)

	headerError := ""

	for i := 0; i < itemsCount; i++ {
		itemes := items[i].(map[string]interface{})
		//ordemData := body["RESULT"].(map[string]interface{})["DATA"].(map[string]interface{})["Ordem"]
		//fmt.Printf("O tipo da variável items é: %T\n", items)
		//log.Info("Executando o item ", itemes["cdMaterial"])
		log.Info("Comeco loop retorno itens:")
		if ordem["cdStatus"].(string) == "3" {
			////////////////////////////////////////////////////////////////////////////
			//Solic APENAS
			if varStuffs["myMap"] == "Solicitacao" {
				parcsz := itemes["Entregas"].([]interface{})
				//items := fileArray["ITENS"].([]interface{})
				parcszCount := len(parcsz)
				var dts []string
				for i := 0; i < parcszCount; i++ {
					parcsze := parcsz[i].(map[string]interface{})
					newDate := parcsze["dtEntregaSolicitada"].(string)
					dts = append(dts, newDate)
				}
				//dates
				var parsedDates []time.Time
				for _, dateStr := range dts {
					parsedDate, err := time.Parse("02012006", dateStr) // Formato ddMMyyyy
					if err != nil {
						log.Error("Erro ao converter data:", err)
					}
					parsedDates = append(parsedDates, parsedDate)
				}

				// Ordenar o slice de time.Time
				sort.Slice(parsedDates, func(i, j int) bool {
					return parsedDates[i].Before(parsedDates[j])
				})

				// Converter de volta para strings no formato original e atualizar o slice original
				for i, date := range parsedDates {
					dts[i] = date.Format("02012006")
				}
				itemes["minDt"] = dts[0]
			}
			//Solic APENAS
			////////////////////////////////////////////////////////////////////////////
			log.Info("Itens checkpoint 1:")
			//log.Info(fmt.Sprintf("Executando o item %s", itemes["cdMaterial"]))
			log.Info(fmt.Sprintf("Executando o item %s", itemes[varStuffs["item"]]))
			itemes["idHeaderTable"] = idHeader
			itemes["cdSequence"] = i + 1

			//rows := db.SqlExec(queryItem, itemes)
			log.Info("Itens checkpoint 2:")
			resultado, errMsg := db.SqlExec(tx, queryItem, itemes)
			fmt.Println("Item - ", itemes[varStuffs["item"]], " ID gerado:", resultado)
			erroLog := ""
			cdStatus := "3"
			if errMsg != "" {
				//erroLog = fmt.Sprintf("Erro no cdMaterial %s - %s", itemes["cdMaterial"], errMsg)
				erroLog = fmt.Sprintf("Erro no %s %s - %s ", varStuffs["item"], itemes[varStuffs["item"]], errMsg)
				fmt.Println(erroLog)
				log.Error(erroLog)
				headerError = headerError + erroLog
				erroLog = errMsg
				fmt.Sprintf("Erro no %s %s - %s", varStuffs["item"], itemes[varStuffs["item"]], errMsg)
				cdStatus = "4"
				tx.Rollback()
			}

			idItem := resultado

			itemes["idItemTable"] = idItem

			////////////////////////////////////
			//Verificando se JSON possui parcelas
			parcs, ok := itemes["Entregas"].([]interface{})
			if ok {
				////////////////////////////////////
				//Chamando função de loop de parcelas
				//////////////////////////
				var campos []string
				for key := range itemes {
					campos = append(campos, key)
				}
				parcsCount := len(parcs)
				for i := 0; i < parcsCount; i++ {
					//itemes := items[i].(map[string]interface{})
					parc := parcs[i].(map[string]interface{})
					for c := 0; c < len(campos); c++ {
						if campos[c] != "Entregas" {
							parc[campos[c]] = itemes[campos[c]]
						}
					}
				}
				//////////////////////////
				parcError := parcInsert(tx, parcs, data, ordem, idItem, idHeader, itemes["cdMaterial"].(string), varStuffs)
				headerError = headerError + parcError
			} else {
				itens := map[string]interface{}{
					varStuffs["item"]: "",
					"cdStatus":        cdStatus,
					"dsLogError":      erroLog,
				}

				itens[varStuffs["item"]] = itemes[varStuffs["item"]]

				ordem["Itens"] = append(ordem["Itens"].([]interface{}), itens)

				log.Info(ordem)

				log.Info("Não tem query para esse layout: ")
				fmt.Println("Não tem query para esse layout: ")
			}
		} else {
			log.Info("Passou aqui no erro cabeças:")

			if varStuffs["myMap"] == "Ordem" {
				parcs := itemes["Entregas"].([]interface{})
				parcsCount := len(parcs)
				for i := 0; i < parcsCount; i++ {
					parc := parcs[i].(map[string]interface{})
					itens := map[string]interface{}{
						varStuffs["parc"]: "",
						"cdStatus":        4,
						"dsLogError":      "Erro no cabeçalho",
					}
					itens[varStuffs["parc"]] = parc[varStuffs["parc"]]
					ordem["Itens"] = append(ordem["Itens"].([]interface{}), itens)
				}
			} else {
				itens := map[string]interface{}{
					varStuffs["item"]: "",
					"cdStatus":        4,
					"dsLogError":      "Erro no cabeçalho",
				}

				itens[varStuffs["item"]] = itemes[varStuffs["item"]]

				ordem["Itens"] = append(ordem["Itens"].([]interface{}), itens)

			}
			/*
				itens := map[string]interface{}{
					iJsonName:    "",
					"cdStatus":   4,
					"dsLogError": "Erro no cabeçalho",
				}

				itens[varStuffs["item"]] = itemes[varStuffs["item"]]

				ordem["Itens"] = append(ordem["Itens"].([]interface{}), itens)
			*/
			log.Info(ordem)

			log.Info("Passei aqui irmao " + itemes["cdMaterial"].(string))
			fmt.Println("Não tem query para esse layout: ")
		}
	}
	return headerError
	//log.Info("Terminado esse item")
}

func parcInsert(tx *sqlx.Tx, parcs []interface{}, data map[string]interface{}, ordem map[string]interface{}, idItem int, idHeader int, codItem string, varStuffs map[string]string) string {
	queryParc := data["PARC"].(string)
	parcsCount := len(parcs)
	log.Info("Total de parc: ", parcsCount)
	itemError := ""

	for i := 0; i < parcsCount; i++ {
		parc := parcs[i].(map[string]interface{})
		log.Info(fmt.Sprintf("Executando a parc %s", varStuffs["parc"]))

		parc["idItemTable"] = idItem
		parc["idHeaderTable"] = idHeader
		parc["cdMaterial"] = codItem
		parc["cdSequence"] = i + 1

		log.Info("parcInsert ITEM - ", parc["idItemTable"])
		log.Info("parcInsert HEAD - ", parc["idHeaderTable"])
		log.Info("parcInsert MATERIAL - ", parc["idHeaderTable"])

		erroLog := ""
		cdStatus := "3"
		//rows := db.SqlExec(queryParc, parc)
		resultado, errMsg := db.SqlExec(tx, queryParc, parc)
		fmt.Println("Parc - ", parc[varStuffs["parc"]], " ID gerado:", resultado)
		if errMsg != "" {
			//erroLog = fmt.Sprintf("Erro no IdPoParc %s - %s", parc["IdPoParc"], errMsg)
			erroLog = fmt.Sprintf("Erro no %s %s - %s ", varStuffs["parc"], parc[varStuffs["parc"]], errMsg)
			fmt.Println(erroLog)
			log.Error(erroLog)
			itemError = itemError + erroLog
			erroLog = errMsg
			cdStatus = "4"
			tx.Rollback()
		}

		parcResultado := resultado
		parc["idParcTable"] = parcResultado
		// Print the results
		//fmt.Printf("Parc Resultado: %d\n", parcResultado)
		//log.Info("ENTREGAS A INSERIR: ")
		//log.Info(parc)

		// Preenchendo os valores posteriormente
		itens := map[string]interface{}{
			varStuffs["parc"]: "",
			"cdStatus":        cdStatus,
			"dsLogError":      erroLog,
		}

		itens[varStuffs["parc"]] = parc[varStuffs["parc"]]
		ordem["Itens"] = append(ordem["Itens"].([]interface{}), itens)

	}
	return itemError
}
