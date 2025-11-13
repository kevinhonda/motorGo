package actions

import "os"

type getNameAction struct{}

func NewGetNameAction() Action {
	return &getNameAction{}
}

func (a *getNameAction) Validate(input []byte) error {
	return nil
}

func (a *getNameAction) Execute() (interface{}, error) {
	return os.Getenv("COMPANY_NAME"), nil
}
