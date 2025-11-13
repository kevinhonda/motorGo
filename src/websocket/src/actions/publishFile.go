package actions

import (
	"errors"

	"github.com/kafy11/gosocket/log"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/file"
	"motorv2/src/websocket/pkg/file"
)

type publishFileAction struct {
	Params *publishFileActionJson

	callback func(filePath, content string)
}

type publishFileActionJson struct {
	FilePath    string `json:"file"`
	FileContent string `json:"content"`
}

func NewPublishFileAction(callback func(filePath, content string)) Action {
	return &JsonActionDecorator{&publishFileAction{
		callback: callback,
	}}
}

func (a *publishFileAction) getJsonParams() interface{} {
	a.Params = &publishFileActionJson{}
	return a.Params
}

func (a *publishFileAction) Validate() error {
	if a.Params.FilePath == "" {
		return errors.New("Missing a file path")
	}

	if a.Params.FileContent == "" {
		return errors.New("Missing a file content")
	}

	return nil
}

func (a *publishFileAction) Execute() (interface{}, error) {
	err := file.Write(a.Params.FilePath, a.Params.FileContent)
	if err != nil {
		return nil, err
	}

	abs, err := file.GetAbs(a.Params.FilePath)
	if err == nil {
		log.Info("Editado o arquivo", abs)
		go a.callback(abs, a.Params.FileContent)
	}

	return "Arquivo publicado", nil
}
