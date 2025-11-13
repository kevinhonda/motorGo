package basicAuth

import (
	"encoding/base64"
	"fmt"
)

func GetBase64(user, password string) string {
	auth := fmt.Sprintf("%s:%s", user, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
