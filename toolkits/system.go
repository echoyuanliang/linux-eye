package toolkits

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var log = logging.MustGetLogger("linux-eye")

type UpTime struct {
	Days    uint64
	Hours   uint8
	Minutes uint8
}

type SystemInfo struct {
	HostName string  `json:"host_name"`
	Platform string  `json:"platform"`
	Os       string  `json:"os"`
	Kernel   string  `json:"kernel"`
	UpTime   *UpTime `json:"up_time"`
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
		log.Error(errMsg)
		return "", errors.New(errMsg)
	}

	return out.String(), nil
}

func SystemUptime() (*UpTime, error) {
	bs, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		errMsg := fmt.Sprintf("read /proc/uptime failed: %v", err)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	content := strings.TrimSpace(string(bs))
	fields := strings.Fields(content)
	if len(fields) < 2 {
		errMsg := "/proc/uptime format not supported"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	secStr := fields[0]
	var secF float64
	secF, err = strconv.ParseFloat(secStr, 64)
	if err != nil {
		errMsg := "/proc/uptime format not supported"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	minTotal := secF / 60.0
	hourTotal := minTotal / 60.0

	days := int64(hourTotal / 24.0)
	hours := int64(hourTotal) - days*24
	minutes := int64(minTotal) - (days * 60 * 24) - (hours * 60)

	return &UpTime{Days: uint64(days), Hours: uint8(hours), Minutes: uint8(minutes)}, nil
}

func GetSystemInfo() (*SystemInfo, error) {

	out, err := getInfo()
	if err != nil || strings.Index(out, "broken pipe") != -1 {
		errMsg := fmt.Sprintf("get system info failed, error: %v, out: %s", err, out)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	osInfo := strings.Split(strings.Replace(out, "\n", " ", -1), " ")
	systemInfo := &SystemInfo{Kernel: strings.TrimSpace(osInfo[1]),
		Platform: strings.TrimSpace(osInfo[2]),
		Os:       strings.TrimSpace(osInfo[3])}

	hostname, err := os.Hostname()
	if err != nil {
		errMsg := fmt.Sprintf("get hostname failed: %v", err)
		log.Error(errMsg)
		return systemInfo, errors.New(errMsg)
	}

	systemInfo.HostName = hostname

	uptime, err := SystemUptime()
	if err != nil {
		errMsg := fmt.Sprintf("get system uptime failed, error: %v", err)
		log.Error(errMsg)
		return systemInfo, errors.New(errMsg)
	}

	systemInfo.UpTime = uptime
	return systemInfo, nil
}
