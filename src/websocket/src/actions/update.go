package actions

import (
	"errors"

	"log"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/update"
	"motorv2/src/websocket/pkg/update"
)

type updateAction struct {
	Params *updateActionJson
}

type updateActionJson struct {
	Source string `json:"source"`
}

func NewUpdateAction() Action {
	return &JsonActionDecorator{&updateAction{}}
}

func (a *updateAction) getJsonParams() interface{} {
	a.Params = &updateActionJson{}
	return a.Params
}

func (a *updateAction) Validate() error {
	if a.Params.Source == "" {
		return errors.New("Missing a source")
	}

	return nil
}

func (a *updateAction) Execute() (interface{}, error) {
	err := update.Do(a.Params.Source)
	if err != nil {
		return nil, err
	}

	defer log.Fatal("Stop service")

	return "Atualizado com sucesso!", nil
}
