package toolkits

import (
	"strings"
	"linux-eye/util"
)



func KernelParam ()(map[string]string, error){
	paramStr, err := util.Exec("sysctl", "-a")
	if err != nil{
		return nil, err
	}

	param := make(map[string]string)

	for _,line := range strings.Split(paramStr, "\n"){
		fields := strings.Split(line, "=")

		if len(fields) != 2{
			continue
		}

		key := strings.TrimSpace(fields[0])
		val := strings.TrimSpace(fields[1])
		param[key] = val
	}

	return param, nil
}


