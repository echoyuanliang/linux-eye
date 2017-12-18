package toolkits

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)


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
		return "", fmt.Errorf("exec %s failed: %v", "uname -srio", err)
	}

	return out.String(), nil
}

func SystemUptime() (*UpTime, error) {
	uptimeFile := "/proc/uptime"
	bs, err := ioutil.ReadFile(uptimeFile)
	if err != nil {
		return nil, fmt.Errorf("read %s failed: %v", uptimeFile, err)
	}

	content := strings.TrimSpace(string(bs))
	fields := strings.Fields(content)
	if len(fields) < 2 {
		return nil, fmt.Errorf("/proc/uptime format not supported")
	}

	secStr := fields[0]
	var secF float64
	secF, err = strconv.ParseFloat(secStr, 64)
	if err != nil {
		return nil, fmt.Errorf("/proc/uptime format not supported")
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
		return nil, fmt.Errorf("get system info failed, error: %v, out: %s", err, out)
	}

	osInfo := strings.Split(strings.Replace(out, "\n", " ", -1), " ")
	systemInfo := &SystemInfo{Kernel: strings.TrimSpace(osInfo[1]),
		Platform: strings.TrimSpace(osInfo[2]),
		Os:       strings.TrimSpace(osInfo[3])}

	hostname, err := os.Hostname()
	if err != nil {
		return systemInfo, fmt.Errorf("get hostname failed: %v", err)
	}

	systemInfo.HostName = hostname

	uptime, err := SystemUptime()
	if err != nil {
		return systemInfo, fmt.Errorf("get system uptime failed, error: %v", err)
	}

	systemInfo.UpTime = uptime
	return systemInfo, nil
}
