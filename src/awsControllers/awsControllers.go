package awsControllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	//"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kafy11/gosocket/log"
	"github.com/valyala/fasthttp"

	"motorv2/pkg/ws"
)

func SendParams(params map[string]string, link, username, password string) string {
	url := link + "/Adm_set_env_config"
	log.Info(fmt.Sprintf("%s - Rodando ", url))

	// Dados que você quer enviar no corpo do POST
	log.Info("SendParams - Convertendo JSON")
	jsonData, err := json.Marshal(params)
	if err != nil {
		errStr := fmt.Sprintf("Erro ao converter para JSON: %s", err)
		log.Error("SendParams- ", errStr)
		fmt.Println(errStr)
		return errStr
	}

	// Fazendo a requisição POST
	log.Info("SendParams - Criando requisição POST")
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		errStr := fmt.Sprintf("Erro na criação da requisição: %s", err)
		log.Error("SendParams- ", errStr)
		//fmt.Println(errStr)
		//fmt.Println("Erro na criação da requisição:", err)
		return errStr
	}
	// Adicionando autenticação básica ao cabeçalho
	log.Info("SendParams - Adicionando autenticação básica ao cabeçalho")
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Erro na requisição POST:", err)
		errStr := fmt.Sprintf("Erro na requisição POST: %s", err)
		log.Error("SendParams- ", errStr)
		//fmt.Println(errStr)
		return errStr
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	log.Info("SendParams - Lendo o corpo da resposta")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("Erro ao ler a resposta:", err)
		errStr := fmt.Sprintf("Erro ao ler a resposta: %s", err)
		log.Error("SendParams- ", errStr)
		//fmt.Println(errStr)
		return errStr
	}

	// Verificando o código de status da resposta
	log.Info("SendParams - Verificando o código de status da resposta")
	if resp.StatusCode != http.StatusOK {
		//fmt.Println("Resposta não OK. Código de status:", resp.StatusCode)
		errStr := string(body)
		fmt.Println("Status code: ", resp.StatusCode)
		log.Error("Status code: ", resp.StatusCode)
		fmt.Println(string(body))
		log.Error("SendParams - ", errStr)
		//fmt.Println(errStr)
		return errStr
	}
	log.Info("Status code: ", resp.StatusCode)
	fmt.Println("Status code: ", resp.StatusCode)

	// Exibindo a resposta
	//fmt.Println("Resposta:", string(body))
	response := fmt.Sprintf("%s - Resposta: %s", url, string(body))
	fmt.Println(response)
	log.Info(response)
	return ""
}

func GetInfos(controller string) (map[string]interface{}, error) {
	log.Info("Controller: " + controller)
	// Leia o corpo da resposta da API
	body, err := ws.Get(controller)
	if err != nil {
		log.Error(fmt.Sprintf("GetInfos - %s: %s", controller, err))
		fmt.Println("GetInfos Get:", controller, " - ", err)
		return nil, err
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Acessando query do retorno do serviço
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Estrutura de dados para armazenar o JSON
	var dataBody map[string]interface{}
	// Faz o parsing do JSON
	err = json.Unmarshal(body, &dataBody)
	if err != nil {
		log.Error(fmt.Sprintf("GetInfos - %s: %s", controller, err))
		log.Error(body)
		fmt.Println("GetInfos Json:", controller, " - ", err)
		//return nil, fmt.Errorf("Erro ao fazer o parsing do JSON: %w", err)
		return nil, err
	}
	log.Info("GetInfo - Controller: " + controller + " - Sucesso")
	return dataBody, nil
}

func SendToS3(fileName, preSigneds string) error {
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Envio para o s3
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Abra o arquivo que você deseja enviar
	//var err error
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("SendToS3 :", fileName, " - ", err)
		log.Error(fmt.Sprintf("SendToS3: %s Open File - %s ", fileName, err))
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error(fmt.Sprintf("SendToS3: %s FileInfo - %s ", fileName, err))
		fmt.Println("SendToS3 :", fileName, " - ", err)
		return err
	}

	// Crie uma solicitação HTTP PUT com Fasthttp
	restReq := fasthttp.AcquireRequest()
	restResp := fasthttp.AcquireResponse()
	restReq.Header.SetMethod("PUT")

	url := preSigneds

	// Defina o URL da solicitação
	restReq.SetRequestURI(url)

	// Abra e envie o corpo do arquivo
	restReq.SetBodyStream(file, int(fileInfo.Size()))

	// Crie um cliente HTTP Fasthttp
	restClient := &fasthttp.Client{}

	// Faça a solicitação para fazer o upload do arquivo
	fmt.Println("SendToS3 - Tentando mandar " + fileName)
	if err := restClient.Do(restReq, restResp); err != nil {
		fmt.Println("restClient:", fileName, " - ", err)
		log.Error(fmt.Sprintf("SendToS3: %s restClient - %s ", fileName, err))
		return err
	}

	// Verifique o código de status da resposta
	if restResp.StatusCode() == fasthttp.StatusOK {
		fmt.Println("SendToS3 - Arquivo enviado com sucesso! - " + fileName)
		log.Info("Arquivo enviado com sucesso! - " + fileName)
	} else {
		fmt.Println(fmt.Printf("SendToS3 - Erro: Código de status inesperado: %d\n", restResp.StatusCode()))
		log.Error(fmt.Sprint("SendToS3 - %s : Código de status inesperado:\n", fileName, restResp.StatusCode()))

	}

	//Remover o arquivo
	/**/
	err = os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}
