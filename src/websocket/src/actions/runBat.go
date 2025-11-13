package actions

import (
	"errors"
	"fmt"
	"os"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/cmd"
	"motorv2/src/websocket/pkg/cmd"
)

type runBatAction struct {
	Params *runBatActionJson
}

type runBatActionJson struct {
	Path string `json:"path"`
	Bat  string `json:"bat"`
}

func NewRunBatAction() Action {
	return &JsonActionDecorator{&runBatAction{}}
}

func (a *runBatAction) getJsonParams() interface{} {
	a.Params = &runBatActionJson{}
	return a.Params
}

func (a *runBatAction) Validate() error {
	if a.Params.Path == "" {
		return errors.New("Missing a path")
	}

	if a.Params.Bat == "" {
		return errors.New("Missing a bat file")
	}

	return nil
}

func (a *runBatAction) Execute() (interface{}, error) {
	absolutePath := fmt.Sprintf("%s/%s", a.Params.Path, a.Params.Bat)

	err := cmd.RunBat(os.Getenv("DRIVER_LETTER"), absolutePath)
	if err != nil {
		return nil, err
	}

	return true, nil
}
