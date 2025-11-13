package file

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

func Write(path string, content ...interface{}) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprint(f, content...)
	if err != nil {
		return err
	}

	return nil
}

func WriteBase64(path, content string) error {
	idx := strings.Index(content, ";base64,")
	if idx < 0 {
		return errors.New("Invalid content")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(content[idx+8:]))
	buff := bytes.Buffer{}
	_, err := buff.ReadFrom(reader)
	if err != nil {
		return err
	}

	return Write(path, string(buff.Bytes()))
}
