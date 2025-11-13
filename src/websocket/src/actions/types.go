package actions

import (
	"encoding/json"
	"errors"
)

type Action interface {
	Validate([]byte) error
	Execute() (interface{}, error)
}

type JsonAction interface {
	getJsonParams() interface{}
	Validate() error
	Execute() (interface{}, error)
}

type JsonActionDecorator struct {
	action JsonAction
}

func (d *JsonActionDecorator) Validate(input []byte) error {
	jsonParams := d.action.getJsonParams()
	if err := json.Unmarshal(input, jsonParams); err != nil {
		return errors.New("Invalid parameters")
	}

	return d.action.Validate()
}

func (d *JsonActionDecorator) Execute() (interface{}, error) {
	return d.action.Execute()
}
