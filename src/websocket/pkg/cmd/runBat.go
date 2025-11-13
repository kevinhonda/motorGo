package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func RunBat(driverLetter, path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	return exec.Command("CMD", fmt.Sprintf("/%s", driverLetter), path).Run()
}
