package client

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kafy11/gosocket/log"
)

type WsClient struct {
	dialer  *websocket.Dialer
	headers http.Header
	url     string
	conn    *websocket.Conn
	m       sync.RWMutex
	sendM   sync.Mutex
}

type WsClientParams struct {
	SSL      bool
	URL      string
	Headers  http.Header
	User     string
	Password string
}

type Handler func(string)

func New(params *WsClientParams) (*WsClient, error) {
	var endpoint string
	if params.SSL {
		endpoint = fmt.Sprintf("wss://%s", params.URL)
	} else {
		endpoint = fmt.Sprintf("ws://%s", params.URL)
	}
	if params.User != "" && params.Password != "" {
		auth := fmt.Sprintf("%s:%s", params.User, params.Password)
		auth_encoded := base64.StdEncoding.EncodeToString([]byte(auth))
		endpoint = fmt.Sprintf("%s?id=0&authorization=%s", endpoint, auth_encoded)
	}

	dialer := websocket.DefaultDialer
	if params.SSL {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return &WsClient{
		dialer:  dialer,
		headers: params.Headers,
		url:     endpoint,
	}, nil
}

func (client *WsClient) ListenMessages(messageHandler Handler) error {
	log.Info("ListenMessageiraaaas")
	messageChannel := make(chan string)
	errChannel := make(chan error)
	go readMessages(client, messageChannel, errChannel)

	for {
		select {
		case message := <-messageChannel:
			go messageHandler(message)

		case err := <-errChannel:
			return err
		}
	}
}

// //////////////////////////////////////////////
func readMessages(client *WsClient, incomingMessages chan string, errChannel chan error) {
	for {
		//cria um lock de leitura
		client.m.RLock()

		_, message, err := client.conn.ReadMessage()
		if err != nil {
			errChannel <- err

			//Desbloqueia a o lock de leitura
			client.m.RUnlock()
			return
		}
		incomingMessages <- string(message)

		//Desbloqueia a o lock de leitura
		client.m.RUnlock()
	}
}

// Função auxiliar para truncar strings nos logs
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// //////////////////////////////////////////////
func (client *WsClient) Connect() error {
	//locka para bloquear leituras na variável enquanto estiver tentando conectar
	client.m.Lock()
	defer client.m.Unlock() //defer para desbloqueiar a variável no final da função
	ws, _, err := client.dialer.Dial(client.url, client.headers)
	if err != nil {
		log.Error("CLIENT.GO - Falha ao conectar no websocket: ", err)
		return err
	}

	client.conn = ws
	return nil
}

/*
	func readMessages(client *WsClient, incomingMessages chan string, errChannel chan error) {
		log.Info("Rodando readMessages") //Teste Logs
		for {
			//cria um lock de leitura
			client.m.RLock()

			if client.conn == nil {
				client.m.RUnlock()
				errChannel <- fmt.Errorf("conexão WebSocket está nil")
				log.Error("Conexão WebSocket está nil") //Temp
				return
			}

			_, message, err := client.conn.ReadMessage()
			if err != nil {
				errChannel <- err
				log.Error("Algo deu errado no readMessages") //Temp
				log.Error(err)                               //Temp
				//Desbloqueia a o lock de leitura
				client.m.RUnlock()
				return
			}
			incomingMessages <- string(message)

			//Desbloqueia a o lock de leitura
			client.m.RUnlock()
		}
	}
*/
func (client *WsClient) Send(message interface{}) error {
	//cria um lock de leitura na conexão
	client.m.RLock()
	defer client.m.RUnlock() //defer para desbloquear o lock no final da função

	client.sendM.Lock() //lock para não ter envios paralelos
	defer client.sendM.Unlock()

	err := client.conn.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}

func (client *WsClient) Stop() error {
	return client.conn.Close()
}
