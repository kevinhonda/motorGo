package ws

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"bytes"
	"encoding/base64"

	"github.com/kafy11/gosocket/log"
)

type basicAuth struct {
	user     string
	password string
}

var auth basicAuth
var wsBaseUrl string

func SetAuth(user, password string) {
	auth = basicAuth{user, password}
}

func SetBaseUrl(baseUrl string) {
	wsBaseUrl = baseUrl
}

func Post(controller string, params map[string]string) ([]byte, error) {
	form := url.Values{}

	for key, value := range params {
		form.Add(key, value)
	}

	url := fmt.Sprintf("%s/%s", wsBaseUrl, controller)

	fmt.Println(form.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(auth.user, auth.password)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func Get(controller string) ([]byte, error) {
	// Crie um cliente HTTP personalizado
	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", wsBaseUrl, controller)
	log.Info(url)

	// Crie uma solicitação GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar a solicitação: %w", err)

	}

	// Adicione as credenciais de autenticação básica ao cabeçalho da solicitação
	req.SetBasicAuth(auth.user, auth.password)

	// Realize a solicitação GET com autenticação básica
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Erro ao fazer a solicitação GET: %w", err)
	}
	defer resp.Body.Close()

	// Verifique o código de status da resposta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("A solicitação retornou um código de status: %s", resp.Status)
	}

	// Leia o corpo da resposta da API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler o corpo da resposta: %w", err)
	}

	return body, nil
}

func SendReturn(jsonData []uint8, controller string) {
	// Crie a requisição POST
	url := fmt.Sprintf("%s/%s", wsBaseUrl, controller)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		log.Error("Erro ao criar a requisição:", err)
		return
	}

	// Adicione os cabeçalhos necessários
	req.Header.Set("Content-Type", "application/json")

	// Adicione autenticação básica
	username := auth.user     // Substitua pelo seu nome de usuário
	password := auth.password // Substitua pela sua senha
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)

	// Envie a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição:", err)
		log.Error("Erro ao enviar a requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Leia a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		log.Error("Erro ao ler a resposta:", err)
		return
	}

	// Imprima a resposta
	log.Info(controller)
	log.Info(controller, " - Resposta:", string(body))
	fmt.Println("Resposta:", string(body))
}
