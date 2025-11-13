package message

import (
	"github.com/kafy11/gosocket/log"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/actions"
	"motorv2/src/websocket/src/actions"
)

type Handler struct {
	actions map[string]actions.Action
}

func NewHandler() *Handler {
	return &Handler{
		actions: make(map[string]actions.Action),
	}
}

func (handler *Handler) AddActionHandler(actionType string, action actions.Action) {
	handler.actions[actionType] = action
}

func runAction(action actions.Action, messageReceived []byte) (interface{}, error) {
	log.Info("Action: ", action) //Kevs
	err := action.Validate(messageReceived)
	if err != nil {
		log.Error("Action Validate error:", err) //Kevs
		return nil, err
	}

	if result, err := action.Execute(); err != nil {
		log.Error("Action Exec error:", err) //Kevs
		return nil, err
	} else {
		return result, nil
	}
}

func (handler *Handler) Run(messageText string) *Response {
	log.Info("Texto recebido", messageText)

	msg, err := NewReceivedMessage([]byte(messageText))
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info("Dados da mensagem", msg)

	if msg.Action == "" {
		return nil
	}

	action, found := handler.actions[msg.Action]
	if !found {
		log.Error("WBsocket error:", err) //Kevs
		return msg.GenerateErrorResponse("Ação não encontrada")
	}

	result, err := runAction(action, msg.Params)
	if err != nil {
		log.Error("WBsocket error:", err) //Kevs
		return msg.GenerateErrorResponse(err.Error())
	}

	return msg.GenerateSuccessResponse(result)
}
