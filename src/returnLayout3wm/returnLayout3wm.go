package returnLayout3wm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"motorv2/pkg/db"
	_ "motorv2/pkg/ws"
	"motorv2/src/returnFuncs"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"
)

/*
var layoutReturn string

var insertsQuerys map[string]interface{}

type GetTest struct{}

returnFuncsInstance := returnFuncs.GetTest{}
*/
// GetTest struct com métodos para manipular os dados globais

var (
	layoutReturn  string
	insertsQuerys map[string]interface{}
	jsonFileName  string
	//returnFuncsInstance *GetTest
	//once sync.Once
)
var returnFuncsInstance = returnFuncs.GetTest{}

func setJsonFileName(name string) {
	jsonFileName = name
}

func GTJson(jsonData []uint8, jsonName string) {
	//log.Info(fmt.Sprintf("O tipo da variável JSONNNNNNN é: %T\n", jsonData))
	fileName := fmt.Sprintf("%s_%s_output.json", layoutReturn, jsonName)

	// Criando o arquivo
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()

	// Escrevendo o JSON no arquivo
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao escrever JSON no arquivo:", err)
		return
	}

	fmt.Printf("JSON escrito com sucesso no arquivo %s\n", fileName)
}

func ReturnStuffs(layouts string, querys map[string]interface{}) {
	fmt.Println("Iniciando Retorno do layout - ", layouts)
	layoutReturn = layouts

	insertsQuerys = querys

	//insertsQuerys = querys

	//fmt.Printf("O tipo da variável testo é: %d\n", countLayouts)

	//Função(getReturns) que puxa o retorno
	//Dentro dessa função é chamado o returnFuncs.CreateJson que transforma o retorno em um arquivo JSON
	log.Info("Wms_inv_header - Passando aqui")
	returnFuncs.GetReturns(layoutReturn, "Wms_inv_header")

	log.Info("Wms_inv_header - Passou aqui")
	//Listando nome de todos os JSONs da pasta jsonFiles do Mrp_po
	filesName := returnFuncs.GetFilesName("jsonFiles", layouts)

	//Exibindo no log o nome dos arquivos
	//log.Info("fileNames - ", filesName)

	//Inserindo as paradas no banco
	dbInserts(filesName)
}

func dbInserts(filesName []string) {
	data := insertsQuerys

	var headersMapList []map[string]interface{}
	var itemMapList []map[string]interface{}

	for _, file := range filesName {
		// Inicia a transação
		tx, err := db.OpenTX()
		log.Info("Header - Chamando db.OpenTX2 de ", file)
		if err != nil {
			log.Error("Erro ao iniciar transação:", err)
		}

		headerTemp := headerInsert(tx, file, data["HEADER"].(string))
		//body, err := returnFuncsInstance.GetReturn("Wms_inv_item?id_invoice=NFe3317123456789")
		//Item Stuffs

		log.Info(fmt.Sprintf("Header - %s - SIT_ERROR:%s", file, headerTemp["SIT_ERROR"]))
		//var itemTemp []map[string]interface{}
		if headerTemp["SIT_ERROR"] == "3" {
			log.Info("Header - Chamando itemInsert do ", file)
			itemTemp, err := itemInsert(tx, file, data["ITEM"].(string))
			if len(itemTemp) > 0 || err == "" {
				itemMapList = append(itemMapList, itemTemp)
				log.Info(fmt.Sprintf("ITEM - ID_INVOICE:%s - SIT_ERROR:%s", headerTemp["ID_INVOICE"], itemTemp["SIT_ERROR"]))
				if itemTemp["SIT_ERROR"] != "3" {
					log.Info("ITEM - Commit do ", file)
					headerTemp["SIT_ERROR"] = 4
					headerTemp["DESC_ERROR"] = fmt.Sprintf("Erro no item: %s", itemTemp["DESC_ERROR"])
				}
			} else {
				headerTemp["SIT_ERROR"] = 4
				headerTemp["DESC_ERROR"] = fmt.Sprintf("Erro no item: %s", err)
			}
		}

		headersMapList = append(headersMapList, headerTemp)

		tx.Commit()

	}
	OuterMap := make(map[string]interface{})
	mapItem := make(map[string]interface{})
	for i, item := range itemMapList {
		key := fmt.Sprintf("ITEM%d", i+1)
		mapItem[key] = item
	}
	OuterMap["ITEMS"] = mapItem
	log.Info("JSON ITEM RETURN:")
	log.Info(OuterMap)

	mapHeader := make(map[string]interface{})
	for i, header := range headersMapList {
		key := fmt.Sprintf("Header%d", i+1)
		mapHeader[key] = header
	}

	log.Info("JSON HEADER RETURN:")
	log.Info(mapHeader)
	fmt.Println("Closed")

	// Converta a estrutura para JSON
	/**/
	jsonHeader, err := json.Marshal(mapHeader)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		log.Error("Erro ao converter para JSON:", err)
		return
	}

	jsonItem, err := json.Marshal(OuterMap)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		log.Error("Erro ao converter para JSON:", err)
		return
	}

	//log.Info("Json que vai ser retornado")
	//log.Info(string(jsonData))
	GTJson(jsonItem, "Item")
	GTJson(jsonHeader, "Header")
	//ws.SendReturn(jsonHeader, "Drp_update_id")

}

func itemInsert(tx *sqlx.Tx, fileName string, queryItem string) (map[string]interface{}, string) {
	fileArray := returnFuncs.JsonToArray(fileName)
	log.Info("EXECUTANDO ITEMS ID_CLIENT_HEADER_PK:", fileArray["ID_CLIENT_HEADER_PK"])
	queryPoronoue := queryItem
	cdStatus := "3"
	erroLog := ""

	//Pegando Itens
	itemController := fmt.Sprintf("Wms_inv_item?id_invoice=%s", fileArray["ID_INVOICE"].(string))
	jsonName := fmt.Sprintf("3wmItem_%s", fileArray["ID_INVOICE"].(string))
	setJsonFileName(jsonName)
	getReturn3wm(itemController)
	//fileItemArray := returnFuncs.JsonToArray(fmt.Sprintf("%s.json", jsonName))
	fileItemArray := Json3wmToArray(fmt.Sprintf("%s.json", jsonName))
	log.Info("fileItemArray:")
	log.Info(fileItemArray)
	log.Info("Tamanho ITENS:", len(fileItemArray))

	//Pegando Itens Batches
	itemBatchController := fmt.Sprintf("wms_inv_item_batch?id_invoice=%s", fileArray["ID_INVOICE"].(string))
	jsonName = fmt.Sprintf("3wmItemBatch_%s", fileArray["ID_INVOICE"].(string))
	setJsonFileName(jsonName)
	getReturn3wm(itemBatchController)
	fileItemBatchArray := Json3wmToArray(fmt.Sprintf("%s.json", jsonName))

	if len(fileItemArray) > 0 {
		//fileArray["QTD_ITENS"] = len(fileItemArray)
		for _, item := range fileItemArray {
			tempItem := item.(map[string]interface{})
			tempItem["ID_CLIENT_HEADER_PK"] = fileArray["ID_CLIENT_HEADER_PK"]
			if len(fileItemBatchArray) > 0 {
				for _, itemBatch := range fileItemBatchArray {
					tempItemBatch := itemBatch.(map[string]interface{})
					//log.Info("tempItemBatch:")
					//log.Info(tempItemBatch)
					if tempItem["ID_SEQ_PK"] == tempItemBatch["SEQ_ITEM"] && tempItem["ID_INVOICE_PK"] == tempItemBatch["ID_INVOICE"] {
						tempItem["NUM_BATCH"] = tempItemBatch["NUM_BATCH"]
						tempItem["NUM_BATCHES"] = tempItemBatch["NUM_BATCHES"]
						tempItem["QTY_BATCH"] = tempItemBatch["QTY_BATCH"]
						tempItem["DATE_MANUFACTURE"] = tempItemBatch["DATE_MANUFACTURE"]
						tempItem["DATE_EXPIRATION"] = tempItemBatch["DATE_EXPIRATION"]
					}
				}
			} else {
				log.Error("SEM ITENS BATCHES PARA O ID_INVOICE:", fileArray["ID_INVOICE"])
			}
			log.Info("tempItem:")
			//log.Info(tempItem)
			log.Info(tempItem["COD_UNIT"])
			erroLog := ""
			resultado, errMsg := db.SqlExec(tx, queryPoronoue, tempItem)
			fmt.Println("ITEM - ID gerado:", resultado)
			if errMsg != "" {
				//erroLog = fmt.Sprintf("Erro no nrDocumentoExterno %s - %s", fileArray["nrDocumentoExterno"], errMsg)
				//log.Error("ITEM - ROLLBACK FILE:", fileName)
				erroLog = fmt.Sprintf("Erro no ID_INVOICE %s COD_ITEM: %s SEQ: %s - %s", fileArray["ID_INVOICE"], tempItem["COD_ITEM"], tempItem["ID_SEQ_PK"], errMsg)
				fmt.Println(erroLog)
				log.Error(erroLog)
				log.Error("ITEM - ROLLBACK:", fileArray["ID_INVOICE"])
				erroLog = errMsg
				cdStatus = "4"
				tx.Rollback()
				break
				//return nil
			}
		}
	} else {
		log.Error("SEM ITENS - ROLLBACK:", fileArray["ID_INVOICE"])
		cdStatus = "4"
		erroLog = "Sem itens para inserir"
		tx.Rollback()
		//return nil, fmt.Sprintf("SEM ITENS ROLLBACK:%s", fileArray["ID_INVOICE"])
	}

	invoice := map[string]interface{}{
		"ID_INVOICE": "",
		"COD_ERROR":  "",
		"SIT_ERROR":  cdStatus,
		"DESC_ERROR": erroLog,
	}
	invoice["ID_INVOICE"] = fileArray["ID_INVOICE"].(string)
	invoice["NUM_ORDER"] = fileArray["NUM_ORDER"].(string)

	log.Info("ITEM 3WM - Terminou insert de itens")
	return invoice, erroLog
}

// func headerInsert(data map[string]interface{}, filesName []string, varStuffs map[string]string) {
func headerInsert(tx *sqlx.Tx, fileName string, queryHeader string) map[string]interface{} {
	queryPoronoue := queryHeader
	fileArray := returnFuncs.JsonToArray(fileName)
	controller := fmt.Sprintf("wms_inv_header_payment?id_invoice=%s", fileArray["ID_INVOICE"].(string))
	jsonName := fmt.Sprintf("3wmHeaderPayment_%s", fileArray["ID_INVOICE"].(string))

	setJsonFileName(jsonName)
	getReturn3wm(controller)
	file3wmHeaderPayment := Json3wmToArray(fmt.Sprintf("%s.json", jsonName))
	if len(file3wmHeaderPayment) > 0 {
		tempHeaderPayment := file3wmHeaderPayment[0]
		tempHeaderPaymentColumns := tempHeaderPayment.(map[string]interface{})
		fileArray["NUM_PAYMENT_PK"] = tempHeaderPaymentColumns["NUM_PAYMENT_PK"]
		fileArray["DATE_EXPIRE"] = tempHeaderPaymentColumns["DATE_EXPIRE"]
		fileArray["VAL_PAYMENT"] = tempHeaderPaymentColumns["VAL_PAYMENT"]
	}

	erroLog := ""
	cdStatus := "3"

	resultado, errMsg := db.SqlExec(tx, queryPoronoue, fileArray)
	fileArray["ID_CLIENT_HEADER_PK"] = resultado
	fmt.Println("Header - ID gerado:", resultado)
	//log.Info(fmt.Sprintf("Executando o header %s", fileArray["nrDocumentoExterno"]))
	log.Info("Header - ID gerado:", resultado)
	log.Info(fmt.Sprintf("Executando o header %s", fileArray["ID_INVOICE"]))

	if errMsg != "" {
		//erroLog = fmt.Sprintf("Erro no nrDocumentoExterno %s - %s", fileArray["nrDocumentoExterno"], errMsg)
		log.Error("HEADER ROLLBACK FILE:", fileName)
		erroLog = fmt.Sprintf("Erro no ID_INVOICE %s - %s", fileArray["ID_INVOICE"], errMsg)
		fmt.Println(erroLog)
		log.Error(erroLog)
		log.Error("HEADER ROLLBACK:", fileArray["ID_INVOICE"])
		erroLog = errMsg
		cdStatus = "4"
		tx.Rollback()
		//return nil
	}

	log.Info("HEADER fileArray ID_CLIENT_HEADER_PK:", fileArray["ID_CLIENT_HEADER_PK"])
	//myMap := make(map[string]interface{})
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
	err = ioutil.WriteFile("jsonFiles/"+fileName, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		log.Error("Erro ao escrever no arquivo:", err)
		//return
	}

	fmt.Printf("Arquivo JSON %s atualizado com sucesso.", fileName)

	ordem := map[string]interface{}{
		"ID_INVOICE": "",
		"NUM_ORDER":  "",
		"SIT_ERROR":  cdStatus,
		"DESC_ERROR": erroLog,
	}
	ordem["ID_INVOICE"] = fileArray["ID_INVOICE"].(string)
	ordem["NUM_ORDER"] = fileArray["NUM_ORDER"].(string)
	//myMap["HEADER"] = ordem

	return ordem

}

func getReturn3wm(controller string) {
	body, err := returnFuncsInstance.GetReturn(controller)
	if err != nil {
		log.Error(fmt.Sprint("getReturn3wm Error:", err))
		//fmt.Println(fmt.Sprint("GetReturn Error:", err))
		//return fmt.Sprint("GetReturn Error:", err)
	}

	fmt.Println("O nome do layout é - ", layoutReturn)
	//jsonsBody := body["RESULT"].(map[string]interface{})["DATA"].(map[string]interface{})[returnLayoutName].([]interface{})
	var jsonsBody ([]interface{})
	if result, ok := body["RESULT"].(map[string]interface{}); ok {
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
					log.Error("Erro ao criar a pasta:", err)
					//return fmt.Sprint("Erro ao criar a pasta:", err)
				}
			} else {
				fmt.Println("Directory exists.")
			}
			//Criando JSONs na pasta
			//returnFuncs.CreateJson("3mw_Item", folderName, jsonsBody)
			////////////////////////////////////////////////////////////////////
			//fileName := fmt.Sprintf("%s/%s", folderName, layoutReturn)
			//fileName := fileArray["ID_INVOICE"].(string)
			CreateJson(folderName, jsonsBody)
		} else {
			log.Error(fmt.Sprintf("%s - A chave 'DATA' não existe.", controller))
			//return fmt.Sprintf("%s - A chave 'DATA' não existe.", controller)
		}

	} else {
		log.Error(fmt.Sprintf("%s - A chave 'RESULT' não existe ou não é do tipo esperado.", controller))
		//return fmt.Sprintf("%s - A chave 'RESULT' não existe ou não é do tipo esperado.", itemController)
	}
}

func CreateJson(folderName string, jsonsBody []interface{}) {

	var headersMapList []interface{}
	for i := 0; i < len(jsonsBody); i++ {
		headersMapList = append(headersMapList, jsonsBody[i])
	}

	file, err := os.Create(folderName + "/" + jsonFileName + ".json")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo JSON:", err)
		log.Error("Erro ao criar o arquivo JSON:", err)
	}
	defer file.Close()

	// Criar um codificador JSON
	encoder := json.NewEncoder(file)

	// Codificar os dados para JSON e escrever no arquivo
	//err = encoder.Encode(jsonsBody)
	err = encoder.Encode(headersMapList)
	if err != nil {
		fmt.Println("Erro ao codificar dados para JSON:", err)
		log.Error("Erro ao codificar dados para JSON:", err)
	}

}

/**/
func Json3wmToArray(jsonName string) []interface{} {
	file, err := os.Open("jsonFiles/" + jsonName)
	log.Info("JsonToArray - Abrindo arquivo:", "jsonFiles/"+jsonName)
	log.Info(file)
	if err != nil {
		fmt.Println("Error:", err)
		log.Error("JsonToArray - Error:", err)
	}
	defer file.Close()

	// Decode JSON from the file
	var data []interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error:", err)
		log.Error("JsonToArray - Error:", err)
	}

	return data
}
