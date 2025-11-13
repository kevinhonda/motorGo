package serverActions

import (
	"github.com/kafy11/gosocket/log"
	//"github.com/kafy11/gowsclient/client"

	"motorv2/src/websocket/src/client"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/message"
	"motorv2/src/websocket/src/message"
)

var wsClient *client.WsClient

func SetWsClient(client *client.WsClient) {
	wsClient = client
}

func CallServer(params map[string]interface{}) {
	msg := message.NewMessageToServer(params)
	log.Info("Enviando mensagem", msg, wsClient)

	err := wsClient.Send(msg)

	if err != nil {
		log.Error("Falha ao enviar mensagem para o server", err)
	}
}
