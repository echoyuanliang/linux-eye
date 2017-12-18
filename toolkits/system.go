package toolkits

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"os/exec"
	"strings"
)

var log = logging.MustGetLogger("linux-eye")

type SystemInfo struct {
	HostName string `json:"host_name"`
	Platform string `json:"platform"`
	Os       string `json:"os"`
	Kernel   string `json:"kernel"`
}

func getInfo() (string, error) {
	cmd := exec.Command("uname", "-srio")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("exec %s failed: %v", "uname -srio", err)
		log.Errorf(errMsg)
		return "", errors.New(errMsg)
	}

	return out.String(), nil
}

func GetSystemInfo() (*SystemInfo, error) {

	out, err := getInfo()
	if err != nil || strings.Index(out, "broken pipe") != -1 {
		errMsg := fmt.Sprintf("get system info failed, error: %v, out: %s", err, out)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}

	osInfo := strings.Split(strings.Replace(out, "\n", " ", -1), " ")
	systemInfo := &SystemInfo{Kernel: strings.TrimSpace(osInfo[1]),
		Platform: strings.TrimSpace(osInfo[2]),
		Os:       strings.TrimSpace(osInfo[3])}

	hostname, err := os.Hostname()
	if err != nil {
		errMsg := fmt.Sprintf("get hostname failed: %v", err)
		log.Errorf(errMsg)
		return systemInfo, errors.New(errMsg)
	}

	systemInfo.HostName = hostname
	return systemInfo, nil

}
