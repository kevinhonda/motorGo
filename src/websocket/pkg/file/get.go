package file

import (
	"io/ioutil"
)

func Get(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
