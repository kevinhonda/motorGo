package actions

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/file"
	"motorv2/src/websocket/pkg/file"
)

type publishZipAction struct {
	Params *publishZipActionJson
}

type publishZipActionJson struct {
	Path       string `json:"path"`
	ZipContent string `json:"content"`
}

func NewPublishZipAction() Action {
	return &JsonActionDecorator{&publishZipAction{}}
}

func (a *publishZipAction) getJsonParams() interface{} {
	a.Params = &publishZipActionJson{}
	return a.Params
}

func (a *publishZipAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}

	if a.Params.ZipContent == "" {
		return errors.New("Missing a zip content")
	}

	return nil
}

func (a *publishZipAction) Execute() (interface{}, error) {
	absolutePath, err := filepath.Abs(a.Params.Path)
	if err != nil {
		return nil, err
	}

	zipPath := fmt.Sprintf("%s%spublish.zip", absolutePath, string(os.PathSeparator))

	err = file.WriteBase64(zipPath, a.Params.ZipContent)
	if err != nil {
		return nil, err
	}

	err = file.Unzip(zipPath, absolutePath)
	if err != nil {
		return nil, err
	}

	err = file.Delete(zipPath)
	if err != nil {
		return nil, err
	}

	return 1, nil
}
