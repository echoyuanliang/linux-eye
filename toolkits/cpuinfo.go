package toolkits

import (
	"bufio"
	"bytes"
	"fmt"
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
	cpuinfoFile := "/proc/cpuinfo"
	bs, err := ioutil.ReadFile(cpuinfoFile)
	if err != nil {
		return nil, fmt.Errorf("read %s failed: %v", cpuinfoFile, err)
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))
	cpuInfo := &CpuInfo{Num: runtime.NumCPU()}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return cpuInfo, fmt.Errorf("read %s buffer failed: %v", cpuinfoFile, err)
		}

		arr := strings.Split(line, ":")
		if len(arr) != 2 {
			continue
		}

		itemName := strings.TrimSpace(arr[0])
		if itemName == "cpu MHz" {

			mHz, err := strconv.ParseFloat(strings.TrimSpace(arr[1]), 32)
			if err != nil {
				return cpuInfo, fmt.Errorf("unsupport %s format: %v", cpuinfoFile, err)
			}

			cpuInfo.MHz = float32(mHz)
		}

		if itemName == "cache size" {
			cpuInfo.CacheSize = strings.TrimSpace(arr[1])
		}
	}

	return cpuInfo, nil

}

func Get()  {
	
}