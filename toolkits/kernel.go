package toolkits

import (
	"strings"
	"linux-eye/util"
	"strconv"
)


func parseParamVal(val string)interface{}{
	if v, err := strconv.ParseInt(val, 10, 64); err == nil{
		return v
	}else if v, err := strconv.ParseFloat(val, 64); err == nil{
		return v
	}

	return val
}



func KernelParam ()(map[string]interface{}, error){
	paramStr, err := util.Exec("sysctl", "-a")
	if err != nil{
		return nil, err
	}

	param := make(map[string]interface{})

	for _,line := range strings.Split(paramStr, "\n"){
		fields := strings.Split(line, "=")

		if len(fields) != 2{
			continue
		}

		key := strings.TrimSpace(fields[0])
		val := parseParamVal(strings.TrimSpace(fields[1]))
		param[key] = val
	}

	return param, nil
}


