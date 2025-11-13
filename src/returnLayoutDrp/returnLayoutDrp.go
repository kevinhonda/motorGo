package returnLayoutDrp

import (
	"encoding/json"
	"fmt"

	"motorv2/pkg/ws"
	"motorv2/src/returnFuncs"
	"os"

	//"strings"

	"github.com/kafy11/gosocket/log"
)

var layoutReturn string

var insertsQuerys map[string]interface{}

func GTJson(jsonData []uint8) {
	//log.Info(fmt.Sprintf("O tipo da variável JSONNNNNNN é: %T\n", jsonData))
	fileName := fmt.Sprintf("%s_output.json", layoutReturn)

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

	//fmt.Printf("O tipo da variável testo é: %d\n", countLayouts)

	//Função(getReturns) que puxa o retorno
	//Dentro dessa função é chamado o returnFuncs.CreateJson que transforma o retorno em um arquivo JSON
	log.Info("Drp_get_id - Passando aqui")
	returnFuncs.GetReturns(layoutReturn, "Drp_get_id")
	log.Info("Drp_get_id - Passou aqui")
	//Listando nome de todos os JSONs da pasta jsonFiles do Mrp_po
	filesName := returnFuncs.GetFilesName("jsonFiles", layouts)
	//Exibindo no log o nome dos arquivos
	//log.Info("fileNames - ", filesName)

	//Inserindo as paradas no banco
	dbInserts(filesName)
}

func dbInserts(filesName []string) {

	//////////////////////////////
	//Puxar scrip do dynamo
	//////////////////////////////

	data := insertsQuerys
	//dataCount := 3
	//fmt.Println("Total dessa fita ai:", dataCount)

	///////////////////////////
	///////////////////////////
	/////Querys de teste///////
	///////////////////////////
	///////////////////////////
	outputReturn := map[string]string{
		"myMap":     "Transferencia",
		"fileArray": "nrDocumentoExterno",
		"header":    "cdTransferencia",
		"item":      "cdMaterial",
		"parc":      "cdMaterial",
	}
	// Preenchendo os valores posteriormente
	// Inicializando um mapa vazio
	//myMap := returnFuncs.HeaderInsert(data, filesName, outputReturn)
	//log.Info(outputReturn["myMap"] + " querys:")
	//log.Info(data)
	returnInsertsInstance := returnFuncs.ReturnInserts{}

	myMap := returnInsertsInstance.HeaderInsert(data, filesName, outputReturn)

	fmt.Println("Closed")

	// Converta a estrutura para JSON
	jsonData, err := json.Marshal(myMap)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		log.Error("Erro ao converter para JSON:", err)
		return
	}

	//log.Info("Json que vai ser retornado")
	//log.Info(string(jsonData))

	GTJson(jsonData)
	ws.SendReturn(jsonData, "Drp_update_id")

}
