package toolkits

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

type CpuInfo struct {
	Num       int     `json:"num"`
	MHz       float32 `json:"m_hz"`
	CacheSize string  `json:"cache_size"`
}

func GetCpuInfo() (*CpuInfo, error) {
	f := "/proc/cpuinfo"

	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))
	cpuInfo := &CpuInfo{Num: runtime.NumCPU()}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return cpuInfo, err
		}

		arr := strings.Split(line, ":")
		if len(arr) != 2 {
			continue
		}

		itemName := strings.TrimSpace(arr[0])
		if itemName == "cpu MHz" {

			mHz, err := strconv.ParseFloat(arr[1], 32)
			if err != nil {
				return cpuInfo, err
			}

			cpuInfo.MHz = float32(mHz)
		}

		if itemName == "cache size" {
			cpuInfo.CacheSize = strings.TrimSpace(arr[1])
		}
	}

	return cpuInfo, nil

}
