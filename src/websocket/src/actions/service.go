package actions

import (
	"errors"

	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/msServices"
	"motorv2/src/websocket/pkg/msServices"
)

type serviceCommandValidation struct {
	validateName bool
}

var serviceCommandOptions map[string]serviceCommandValidation

type serviceAction struct {
	Params *serviceActionJson
}

type serviceActionJson struct {
	Command string `json:"command"`
	Name    string `json:"name"`
	Like    string `json:"like"`
}

func init() {
	serviceCommandOptions = map[string]serviceCommandValidation{
		"started":  {validateName: false},
		"stopped":  {validateName: false},
		"listlike": {validateName: false},
		"stop":     {validateName: true},
		"restart":  {validateName: true},
		"start":    {validateName: true},
	}
}

func NewServiceAction() Action {
	return &JsonActionDecorator{&serviceAction{}}
}

func (a *serviceAction) getJsonParams() interface{} {
	a.Params = &serviceActionJson{}
	return a.Params
}

func (a *serviceAction) Validate() error {
	if a.Params.Command == "" {
		return errors.New("Missing a command")
	}

	command, foundCommand := serviceCommandOptions[a.Params.Command]
	if !foundCommand {
		return errors.New("Invalid service command")
	}

	if command.validateName && a.Params.Name == "" {
		return errors.New("Missing a service name")
	}

	return nil
}

func (a *serviceAction) Execute() (interface{}, error) {
	switch a.Params.Command {
	case "started":
		return msServices.GetStarted()
	case "stopped":
		return msServices.GetStopped()
	case "listlike":
		return msServices.GetLike(a.Params.Like)
	case "stop":
		return "", msServices.Stop(a.Params.Name)
	case "start":
		return "", msServices.Start(a.Params.Name)
	case "restart":
		return "", msServices.Restart(a.Params.Name)
	}
	return nil, errors.New("Invalid service command")
}
