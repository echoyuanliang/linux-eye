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
		cpuBytes, err := json.MarshalIndent(cpuInfo, "", "    ")
		if err != nil{
			log.Errorf("json marshal cpuInfo failed: %v", err)
		}else{
			fmt.Println(string(cpuBytes))
		}
	}


	sysInfo, err := toolkits.GetSystemInfo()

	if err != nil{
		log.Errorf("get sys info failed: %v", err)
	}else{
		sysBytes, err := json.MarshalIndent(sysInfo, "", "    ")
		if err != nil{
			log.Errorf("json marshal sysInfo failed: %v", err)
		}else{
			fmt.Println(string(sysBytes))
		}
	}

	fsInfo, err := toolkits.ListDeviceUsage()
	if err != nil{
		log.Errorf("get fs info failed: %v", err)
	}else{
		sysBytes, err := json.MarshalIndent(fsInfo, "", "    ")
		if err != nil{
			log.Errorf("json marshal fsInfo failed: %v", err)
		}else{
			fmt.Println(string(sysBytes))
		}
	}



	procStat, err := toolkits.CurrentProcStat()
	if err != nil{
		log.Errorf("get proc stat failed: %v", err)
	}else {
		procBytes, err := json.MarshalIndent(procStat, "", "    ")
		if err != nil{
			log.Errorf("json marshal procStat failed: %v", err)
		}else{
			fmt.Println(string(procBytes))
		}
	}

}