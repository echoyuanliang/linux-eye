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
	}else {
		cpuStr, err := json.Marshal(cpuInfo)
		if err != nil{
			log.Errorf("json marshal cpuInfo failed: %v", err)
		}else{
			fmt.Println(cpuStr)
		}
	}


	sysInfo, err := toolkits.GetSystemInfo()

	if err != nil{
		log.Errorf("get sys info failed: %v", err)
	}else{
		sysStr, err := json.Marshal(sysInfo)
		if err != nil{
			log.Errorf("json marshal sysInfo failed: %v", err)
		}else{
			fmt.Println(sysStr)
		}
	}

	fsInfo, err := toolkits.ListDeviceUsage()
	if err != nil{
		log.Errorf("get fs info failed: %v", err)
	}else{
		sysStr, err := json.Marshal(fsInfo)
		if err != nil{
			log.Errorf("json marshal fsInfo failed: %v", err)
		}else{
			fmt.Println(sysStr)
		}
	}



	procStat, err := toolkits.CurrentProcStat()
	if err != nil{
		log.Errorf("get proc stat failed: %v", err)
	}else {
		procStr, err := json.Marshal(procStat)
		if err != nil{
			log.Errorf("json marshal procStat failed: %v", err)
		}else{
			fmt.Println(procStr)
		}

		fmt.Println(procStr)
	}

}