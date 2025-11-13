package actions

import (
	"errors"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/file"
	"motorv2/src/websocket/pkg/file"
)

type deleteFileAction struct {
	Params *deleteFileActionJson
}

type deleteFileActionJson struct {
	Path string `json:"path"`
}

func NewDeleteFileAction() Action {
	return &JsonActionDecorator{&deleteFileAction{}}
}

func (a *deleteFileAction) getJsonParams() interface{} {
	a.Params = &deleteFileActionJson{}
	return a.Params
}

func (a *deleteFileAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}
	return nil
}

func (a *deleteFileAction) Execute() (interface{}, error) {
	err := file.Delete(a.Params.Path)

	if err != nil {
		return nil, err
	}

	return "Deletado arquivo com sucesso", nil
}
