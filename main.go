package main

import (
	"linux-eye/toolkits"
	"fmt"
	"github.com/op/go-logging"
	"encoding/json"
)

var log = logging.MustGetLogger("linux-eye")


func main()  {
	cpuInfo, err := toolkits.GetCpuInfo()
	if err != nil{
		log.Errorf("get cpu info failed: %v", err)
	}

	fmt.Print(json.Marshal(cpuInfo))

	sysInfo, err := toolkits.GetSystemInfo()

	if err != nil{
		log.Errorf("get sys info failed: %v", err)
	}
	fmt.Print(json.Marshal(sysInfo))

	fsInfo, err := toolkits.ListDeviceUsage()
	if err != nil{
		log.Errorf("get fs info failed: %v", err)
	}
	fmt.Print(json.Marshal(fsInfo))

	procStat, err := toolkits.CurrentProcStat()
	if err != nil{
		log.Errorf("get proc stat failed: %v", err)
	}
	fmt.Print(json.Marshal(procStat))

}