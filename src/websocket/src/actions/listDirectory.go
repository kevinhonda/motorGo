package actions

import (
	"errors"
	"io/ioutil"
)

type listDirectoryAction struct {
	Params *listDirectoryActionJson
}

type listDirectoryActionJson struct {
	Path string `json:"directory"`
}

func NewListDirectoryAction() Action {
	return &JsonActionDecorator{&listDirectoryAction{}}
}

func (a *listDirectoryAction) getJsonParams() interface{} {
	a.Params = &listDirectoryActionJson{}
	return a.Params
}

func (a *listDirectoryAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}
	return nil
}

func (a *listDirectoryAction) Execute() (interface{}, error) {
	files, err := ioutil.ReadDir(a.Params.Path)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, file := range files {
		result = append(result, file.Name())
	}

	return result, nil
}
