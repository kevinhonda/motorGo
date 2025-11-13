package serverActions

import "github.com/kafy11/gosocket/log"

const ADD_S3_FILE_KEY = "add_s3_file"

func AddS3File(objectName, body string) {
	log.Info("Adicionando arquivo ao s3", objectName, body)

	CallServer(map[string]interface{}{
		"action":      ADD_S3_FILE_KEY,
		"object_name": objectName,
		"body":        body,
	})
}
