package actions

import (
	"errors"
	"os"
)

type createDirectoryAction struct {
	Params *createDirectoryActionJson
}

type createDirectoryActionJson struct {
	Path string `json:"path"`
}

func NewCreateDirectoryAction() Action {
	return &JsonActionDecorator{&createDirectoryAction{}}
}

func (a *createDirectoryAction) getJsonParams() interface{} {
	a.Params = &createDirectoryActionJson{}
	return a.Params
}

func (a *createDirectoryAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}
	return nil
}

func (a *createDirectoryAction) Execute() (interface{}, error) {
	err := os.Mkdir(a.Params.Path, 0755)

	if err != nil {
		return nil, err
	}

	return "Criado diretório", nil
}
