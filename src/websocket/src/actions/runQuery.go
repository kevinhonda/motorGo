package actions

import (
	"errors"
	"regexp"
	"strings"

	"github.com/kafy11/gosocket/log"
	//"gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/db"
	"motorv2/src/websocket/pkg/db"
)

type runQueryAction struct {
	Params   *runQueryActionJson
	onCreate func(objectName, query string)
}

type runQueryActionJson struct {
	Query string `json:"query"`
}

func NewRunQueryAction(onCreate func(objectName, query string)) Action {
	return &JsonActionDecorator{&runQueryAction{
		onCreate: onCreate,
	}}
}

func (a *runQueryAction) getJsonParams() interface{} {
	a.Params = &runQueryActionJson{}
	return a.Params
}

func (a *runQueryAction) Validate() error {
	if a.Params.Query == "" {
		return errors.New("Missing a query")
	}
	return nil
}

func (a *runQueryAction) Execute() (interface{}, error) {
	if _, err := db.Connection.Exec(a.Params.Query); err != nil {
		return err.Error(), nil
	}

	if a.onCreate != nil {
		log.Info("Checando se foi create")
		objectName := a.checkCreate()
		if objectName != "" {
			objectName = strings.ToUpper(objectName)
			log.Info("Foi create do objeto", objectName)
			go a.onCreate(objectName, a.Params.Query)
		}
	}

	return true, nil
}

func (a *runQueryAction) checkCreate() string {
	r := regexp.MustCompile(`(?i)(?:CREATE|ALTER).*(?:(?P<object_name>(?:PRC_|VW_ERP_)[^"\s]*))`)
	matches := r.FindStringSubmatch(a.Params.Query)

	if len(matches) == 0 {
		return ""
	}
	names := r.SubexpNames()
	log.Info(matches, names)

	for i, name := range names {
		if name == "object_name" {
			return matches[i]
		}
	}
	return ""
}
