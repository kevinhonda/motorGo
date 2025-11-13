package actions

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"

	"github.com/kafy11/gosocket/log"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/db"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/ws"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/zlib"
	"motorv2/src/websocket/pkg/db"
	"motorv2/src/websocket/pkg/wsWebsocket"
	"motorv2/src/websocket/pkg/zlib"
)

type selectAction struct {
	Params *selectActionJson
}

type selectActionJson struct {
	Query        string `json:"query"`
	Queue        string `json:"queue"`
	Table        string `json:"table"`
	Filename     string `json:"filename"`
	ClientSelect int    `json:"clientSelect"`
}

func NewSelectAction() Action {
	return &JsonActionDecorator{&selectAction{}}
}

func (a *selectAction) getJsonParams() interface{} {
	a.Params = &selectActionJson{}
	return a.Params
}

func (a *selectAction) Validate() error {
	if a.Params.Query == "" {
		return errors.New("Missing a query")
	}
	return nil
}

func (a *selectAction) Execute() (interface{}, error) {
	limit := 1000
	if a.Params.ClientSelect == 1 {
		if a.Params.Table != "" {
			limit = -1
		} else {
			limit = 25000
		}
	}

	result, err := db.Connection.MapSelect(a.Params.Query, limit)

	if err != nil {
		return nil, err
	}

	if a.Params.ClientSelect == 1 {
		partSize := 25000
		for len(result) > 0 {
			if len(result) < partSize {
				partSize = len(result)
			}

			data := result[0:partSize]
			result = result[partSize:]

			json, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}

			compressedJson, err := a.compressJson(json)
			if err != nil {
				return nil, err
			}

			_, err = a.sendClientSelectResponse(compressedJson)
			if err != nil {
				return nil, err
			}
		}
		return "", nil
	}
	return result, nil
}

func (a *selectAction) compressJson(json []byte) (string, error) {
	compressedJson, err := zlib.Compress(json, 9)
	if err != nil {
		return "", err
	}
	sEnc := b64.StdEncoding.EncodeToString(compressedJson)
	return sEnc, nil
}

func (a *selectAction) sendClientSelectResponse(response string) ([]byte, error) {
	wsResponse, err := wsWebsocket.Call("adm_wsocket_client_select_return", map[string]string{
		"data":     response,
		"queue":    a.Params.Queue,
		"table":    a.Params.Table,
		"filename": a.Params.Filename,
	})

	if err != nil {
		return nil, err
	}
	log.Info("Resposta do WS", string(wsResponse))

	return wsResponse, nil
}
