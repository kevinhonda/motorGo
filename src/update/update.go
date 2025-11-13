package update

import (
	"net/http"

	"github.com/inconshreveable/go-update"
)

func Do(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return update.Apply(resp.Body, update.Options{})
}
