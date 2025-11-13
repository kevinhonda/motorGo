package actions

import (
	"errors"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/file"
	"motorv2/src/websocket/pkg/file"
)

type getFileAction struct {
	Params *getFileActionJson
}

type getFileActionJson struct {
	Path string `json:"file"`
}

func NewGetFileAction() Action {
	return &JsonActionDecorator{&getFileAction{}}
}

func (a *getFileAction) getJsonParams() interface{} {
	a.Params = &getFileActionJson{}
	return a.Params
}

func (a *getFileAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}
	return nil
}

func (a *getFileAction) Execute() (interface{}, error) {
	content, err := file.Get(a.Params.Path)
	if err != nil {
		return nil, err
	}

	return string(content), nil
}
