package actions

import (
	"errors"
	"os"
)

type deleteDirectoryAction struct {
	Params *deleteDirectoryActionJson
}

type deleteDirectoryActionJson struct {
	Path string `json:"path"`
}

func NewDeleteDirectoryAction() Action {
	return &JsonActionDecorator{&deleteDirectoryAction{}}
}

func (a *deleteDirectoryAction) getJsonParams() interface{} {
	a.Params = &deleteDirectoryActionJson{}
	return a.Params
}

func (a *deleteDirectoryAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}
	return nil
}

func (a *deleteDirectoryAction) Execute() (interface{}, error) {
	err := os.Remove(a.Params.Path)

	if err != nil {
		return nil, err
	}

	return "Deletado o diretório", nil
}
