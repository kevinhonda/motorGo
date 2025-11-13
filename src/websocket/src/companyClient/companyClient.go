package companyClient

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kafy11/gosocket/log"
	//"github.com/kafy11/gowsclient/client"

	"motorv2/src/websocket/src/client"

	//"github.com/kafy11/gowsclient"

	/*"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/basicAuth"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/file"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/actions"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/message"
	"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/serverActions"*/

	"motorv2/src/websocket/pkg/basicAuth"
	"motorv2/src/websocket/pkg/file"
	"motorv2/src/websocket/src/actions"
	"motorv2/src/websocket/src/message"
	"motorv2/src/websocket/src/serverActions"
)

type Client struct {
	client    *client.WsClient
	connected bool
}

type Params struct {
	URL             string
	CompanyName     string
	DB              string
	AuthUser        string
	AuthPass        string
	Version         string
	ControllersPath string
}

var controllersPath string

func New(params *Params) (*Client, error) {
	url := fmt.Sprintf(`%s?id=0&authorization=%s&name=%s&db=%s&version=%s`,
		params.URL,
		basicAuth.GetBase64(params.AuthUser, params.AuthPass),
		url.QueryEscape(params.CompanyName),
		params.DB,
		url.QueryEscape(params.Version))
	log.Info("Websocket URL:", url)
	//cria um client do websocket e connecta

	wsClient, err := client.New(&client.WsClientParams{
		SSL: true,
		URL: url,
	})
	if err != nil {
		log.Error("Websocket Falha ao criar o client:", err)
		return nil, err
	}
	log.Info("Websocket Sucesso ao criar o client")
	serverActions.SetWsClient(wsClient)
	controllersPath = params.ControllersPath

	return &Client{
		client: wsClient,
	}, nil
}

func (c *Client) Run() {
	c.connect()
	log.Info("Websocket Sucesso ao conectar")
	msgHandler := createMessageHandler(c.client)
	c.listenMessages(msgHandler)
}

func (c *Client) Stop() error {
	if c.connected {
		return c.client.Stop()
	}
	return nil
}

func (c *Client) listenMessages(handler *message.Handler) {
	for {
		err := c.client.ListenMessages(func(msgReceived string) {
			response := handler.Run(msgReceived)
			if response == nil {
				log.Info("Websocket handler response nil")
				return
			}

			log.Info("Resposta", response)

			err := c.client.Send(response)
			if err != nil {
				log.Error(err)
			}
		})

		if err != nil {
			log.Error("listenMessages - Falha ao ler mensagem", err) //Kev
			//log.Error("Falha ao ler mensagem", err)
			log.Info("Reconectando")

			c.connect()
		}
	}
}

func (c *Client) connect() {
	c.connected = false
	for !c.connected {
		log.Info("Tentando conectar no websocket")
		err := c.client.Connect()

		if err != nil {
			log.Error("Falha ao conectar:", err)
			time.Sleep(60 * 5 * time.Second)
		} else {
			log.Info("Conectado no websocket") //Kev
			c.connected = true
		}
	}
}

func createMessageHandler(wsClient *client.WsClient) *message.Handler {
	log.Info("Websocket createMessageHandler:", wsClient) //Kev
	messageHandler := message.NewHandler()
	messageHandler.AddActionHandler("call_soap", actions.NewCallSoapAction())
	messageHandler.AddActionHandler("createdir", actions.NewCreateDirectoryAction())
	messageHandler.AddActionHandler("getddl", actions.NewSelectAction())
	messageHandler.AddActionHandler("getfile", actions.NewGetFileAction())
	messageHandler.AddActionHandler("getname", actions.NewGetNameAction())
	messageHandler.AddActionHandler("listdir", actions.NewListDirectoryAction())
	messageHandler.AddActionHandler("publish", actions.NewPublishFileAction(handlePublishFile))
	messageHandler.AddActionHandler("publishzip", actions.NewPublishZipAction())
	messageHandler.AddActionHandler("refresh", actions.NewRefreshAction())
	messageHandler.AddActionHandler("removedir", actions.NewDeleteDirectoryAction())
	messageHandler.AddActionHandler("removefile", actions.NewDeleteFileAction())
	messageHandler.AddActionHandler("runbat", actions.NewRunBatAction())
	messageHandler.AddActionHandler("runquery", actions.NewRunQueryAction(handleCreateDatabaseObject))
	messageHandler.AddActionHandler("select", actions.NewSelectAction())
	messageHandler.AddActionHandler("service", actions.NewServiceAction())
	messageHandler.AddActionHandler("update", actions.NewUpdateAction())

	return messageHandler
}

func (c *Client) Test() {
	c.connect()
	log.Info("Sucesso ao conectar")
}

func handleCreateDatabaseObject(objectName, body string) {
	serverActions.AddS3File(fmt.Sprintf(`DDL/%s.txt`, objectName), body)
}

func handlePublishFile(filePath, content string) {
	controllerAbs, err := file.GetAbs(controllersPath)
	if err != nil {
		log.Error("Falha ao pegar o absolute path do diretório dos controllers", err)
		return
	}

	abs, err := file.GetAbs(filePath)
	if err != nil {
		log.Error("Falha ao pegar o absolute path do arquivo", err)
		return
	}

	if !strings.Contains(abs, controllerAbs) {
		return
	}

	filePath = strings.Replace(filePath, controllerAbs, "", 1)
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	serverActions.AddS3File(fmt.Sprintf(`ws%s`, filePath), content)
}
