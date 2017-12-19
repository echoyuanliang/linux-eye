package main

import (
	"github.com/op/go-logging"
	"linux-eye/toolkits"
	"sync"
	"encoding/json"
	"fmt"
)

var log = logging.MustGetLogger("linux-eye")

type infoMap struct{
	BuiltinMap map[string]interface{}
	Mutex sync.RWMutex
}


func (im *infoMap) Set(key string, val interface{}){
	im.Mutex.Lock()
	im.BuiltinMap[key] = val
	im.Mutex.Unlock()
}


func (im *infoMap) Marshal() (string, error){
	im.Mutex.RLock()
	defer im.Mutex.RUnlock()

	infoBytes, err := json.MarshalIndent(im.BuiltinMap, "", "    ")
	if err != nil{
		return "", fmt.Errorf("marshal infoMap failed: %v", err)
	}else {
		return string(infoBytes), nil
	}
}


func NewInfoMap()(*infoMap){
	return &infoMap{BuiltinMap:make(map[string]interface{})}
}


var infoProcesses = []string{"sys_info", "cpu_info", "kernel_param", "io_stat", "df_stat",
"if_stat", "cpu_stat", "net_stat", "mem_info", "tcp_link"}

func main()  {
	var wg sync.WaitGroup
	wg.Add(len(infoProcesses))

	im := NewInfoMap()

	for _, name := range infoProcesses{
		go func(name string) {

			var info interface{}
			err := fmt.Errorf("%s not defined", name)

			switch name{
				case "sys_info":
					info, err = toolkits.GetSystemInfo()
				case "cpu_info":
					info, err = toolkits.GetCpuInfo()
				case "mem_info":
					info, err = toolkits.MemInfo()
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
				case "net_stat":
					info, err = toolkits.NetStat()
				case "tcp_link":
					info, err = toolkits.TcpLinks()
			}

			if err != nil{
				log.Errorf("get %s failed: %v", name, err)
			}else {
				im.Set(name, info)
			}

			wg.Done()
		}(name)

	}

	wg.Wait()
	infoStr, err := im.Marshal()
	if err != nil{
		log.Fatalf("marshal infoMap failed: %v", err)
	}else {
		fmt.Println(infoStr)
	}
}