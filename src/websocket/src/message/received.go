package message

import "encoding/json"

func NewReceivedMessage(jsonReceived []byte) (*Received, error) {
	msgReceived := &Received{}

	if err := json.Unmarshal(jsonReceived, &msgReceived); err != nil {
		return nil, err
	}

	return msgReceived, nil
}
