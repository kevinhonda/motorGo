package file

import (
	"path/filepath"
)

func GetAbs(path string) (string, error) {
	return filepath.Abs(path)
}
