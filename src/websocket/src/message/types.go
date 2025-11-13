package message

import "encoding/json"

const SEND_MESSAGE_ACTION string = "sendMessage"

type Response struct {
	Action string                 `json:"action"`
	To     string                 `json:"to"`
	Params map[string]interface{} `json:"msg"`
}

type Received struct {
	Action string          `json:"action"`
	From   string          `json:"from"`
	Params json.RawMessage `json:"msg"`
}

type ToServer struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"msg"`
}

func NewMessageToServer(params map[string]interface{}) *ToServer {
	return &ToServer{
		Action: SEND_MESSAGE_ACTION,
		Params: params,
	}
}
