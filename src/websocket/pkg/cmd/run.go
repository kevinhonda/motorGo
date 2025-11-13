package cmd

import (
	"os/exec"
	"strings"
)

func Run(command string, arg ...string) error {
	return exec.Command(command, arg...).Run()
}

func RunOutput(command string, arg ...string) ([]string, error) {
	result, err := exec.Command(command, arg...).Output()

	if err != nil {
		return nil, err
	}
	return strings.Split(string(result), "\r\r\n"), nil
}
