package soap

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func Call(wsdl, xml string) ([]byte, error) {
	payload := strings.NewReader(xml)
	client := &http.Client{}
	req, err := http.NewRequest("POST", wsdl, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/xml")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
