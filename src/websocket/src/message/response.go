package message

import "encoding/json"

func (msg *Received) generateResponse(response interface{}, status int) *Response {
	params := make(map[string]interface{})
	json.Unmarshal(msg.Params, &params)

	params["action"] = msg.Action
	params["response"] = response
	params["status"] = status

	return &Response{
		Action: SEND_MESSAGE_ACTION,
		To:     msg.From,
		Params: params,
	}
}

func (msg *Received) GenerateSuccessResponse(response interface{}) *Response {
	return msg.generateResponse(response, 1)
}

func (msg *Received) GenerateErrorResponse(err string) *Response {
	return msg.generateResponse(err, 0)
}
