package main

import (
	"github.com/op/go-logging"
	"linux-eye/toolkits"
	"sync"
	"encoding/json"
	"fmt"
)

var log = logging.MustGetLogger("linux-eye")
var infoMap = map[string]interface{}{}

var infoProcesses = []string{"sys_info", "cpu_info", "kernel_param", "io_stat", "df_stat",
"if_stat", "cpu_stat", "net_stat", "mem_info"}

func main()  {
	var wg sync.WaitGroup
	wg.Add(len(infoProcesses))

	for _, name := range infoProcesses{
		go func(name string) {

			var info interface{}
			err := fmt.Errorf("%s not defined", name)

			switch name{
				case "sys_info":
					info, err = toolkits.GetSystemInfo()
				case "cpu_info":
					info, err = toolkits.GetCpuInfo()
				case "kernel_param":
					info, err = toolkits.KernelParam()
				case "io_stat":
					info, err = toolkits.ListDiskStats()
				case "df_stat":
					info, err = toolkits.ListDeviceUsage()
				case "if_stat":
					info, err = toolkits.NetIfs()
				case "cpu_stat":
					info, err = toolkits.CurrentProcStat()
				case "mem_info":
					info, err = toolkits.MemInfo()
				case "net_stat":
					info, err = toolkits.NetStat()

			}

			if err != nil{
				log.Errorf("get %s failed: %v", name, err)
			}else {
				infoMap[name] = info
			}

			wg.Done()
		}(name)

	}

	wg.Wait()
	infoBytes, err := json.MarshalIndent(infoMap, "", "    ")
	if err != nil{
		log.Fatalf("marshal infoMap failed: %v", err)
	}else {
		fmt.Println(string(infoBytes))
	}
}