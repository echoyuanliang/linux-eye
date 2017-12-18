package toolkits

import (
	"bufio"
	"bytes"
	"errors"
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

	bs, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		errMsg := fmt.Sprintf("read /proc/cpuinfo failed: %v", err)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))
	cpuInfo := &CpuInfo{Num: runtime.NumCPU()}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			errMsg := fmt.Sprintf("read /proc/cpuinfo buffer failed: %v", err)
			log.Error(errMsg)
			return cpuInfo, errors.New(errMsg)
		}

		arr := strings.Split(line, ":")
		if len(arr) != 2 {
			continue
		}

		itemName := strings.TrimSpace(arr[0])
		if itemName == "cpu MHz" {

			mHz, err := strconv.ParseFloat(strings.TrimSpace(arr[1]), 32)
			if err != nil {
				errMsg := fmt.Sprintf("unsupport /proc/cpuinfo format: %v", err)
				log.Error(errMsg)
				return cpuInfo, errors.New(errMsg)
			}

			cpuInfo.MHz = float32(mHz)
		}

		if itemName == "cache size" {
			cpuInfo.CacheSize = strings.TrimSpace(arr[1])
		}
	}

	return cpuInfo, nil

}
