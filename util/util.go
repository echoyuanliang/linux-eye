package util

import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
)

func Exec(name string, arg... string) (string, error){
	cmd := exec.Command(name, arg...)
	cmd.Stdin = strings.NewReader("cmd input")
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("exec %s failed: %v", "name", arg)
	}

	return stdOut.String(), nil

}