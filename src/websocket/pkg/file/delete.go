package file

import (
	"os"
	"path/filepath"

	"github.com/kafy11/gosocket/log"
)

func Delete(path string) error {
	err := os.Remove(path)

	if err != nil {
		return err
	}

	return nil
}

func DeleteAllWithExtension(directoryPath string, extension string) error {
	d, err := os.Open(directoryPath)
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return err
	}

	log.Info("Encontrado", len(files), "arquivos na pasta", directoryPath)

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == extension {
				filename := file.Name()

				log.Info("Deletando arquivo", filename)
				err = os.Remove(file.Name())
				if err != nil {
					log.Error("Erro deletando arquivo", filename, err)
				}
			}
		}
	}
	return nil
}
