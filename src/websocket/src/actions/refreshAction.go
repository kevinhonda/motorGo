package actions

import "os"

type refreshAction struct{}

func NewRefreshAction() Action {
	return &refreshAction{}
}

func (a *refreshAction) Validate(input []byte) error {
	return nil
}

func (a *refreshAction) Execute() (interface{}, error) {
	os.Exit(0)
	return true, nil
}
