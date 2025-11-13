package wsWebsocket

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type basicAuth struct {
	user     string
	password string
}

var auth basicAuth
var wsBaseUrl string

func SetAuth(user, password string) {
	auth = basicAuth{user, password}
}

func SetBaseUrl(baseUrl string) {
	wsBaseUrl = baseUrl
}

func Call(controller string, params map[string]string) ([]byte, error) {
	form := url.Values{}

	for key, value := range params {
		form.Add(key, value)
	}

	url := fmt.Sprintf("%s/%s", wsBaseUrl, controller)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(auth.user, auth.password)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
