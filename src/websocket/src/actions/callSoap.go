package actions

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/kafy11/gosocket/log"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/soap"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/ws"

	"motorv2/src/websocket/pkg/soap"
	"motorv2/src/websocket/pkg/wsWebsocket"
)

type callSoapAction struct {
	Params *callSoapActionJson
}

type callSoapActionJson struct {
	Type    string `json:"type"`
	Wsdl    string `json:"wsdl"`
	XML     string `json:"xml"`
	QueueID int    `json:"queue"`
}

func NewCallSoapAction() Action {
	return &JsonActionDecorator{&callSoapAction{}}
}

func (a *callSoapAction) getJsonParams() interface{} {
	a.Params = &callSoapActionJson{}
	return a.Params
}

func (a *callSoapAction) Validate() error {
	if a.Params.Type == "" {
		return errors.New("Missing a type")
	}

	if a.Params.Wsdl == "" {
		return errors.New("Missing a wsdl")
	}

	if a.Params.XML == "" {
		return errors.New("Missing a xml")
	}

	return nil
}

func (a *callSoapAction) Execute() (interface{}, error) {
	response, err := soap.Call(a.Params.Wsdl, a.Params.XML)
	if err != nil {
		return nil, err
	}
	log.Info("Resposta do SOAP", string(response))

	sentParams, err := json.Marshal(a.Params)
	if err != nil {
		return nil, err
	}

	wsResponse, err := wsWebsocket.Call("Int_bid_ext_xml_integrator", map[string]string{
		"xml":   string(response),
		"sent":  string(sentParams),
		"type":  a.Params.Type,
		"queue": strconv.Itoa(a.Params.QueueID),
	})
	if err != nil {
		return nil, err
	}
	log.Info("Resposta do WS", string(wsResponse))

	return "1", nil
}
