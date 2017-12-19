package toolkits

import (
	"io/ioutil"
	"fmt"
	"bufio"
	"bytes"
	"io"
	"strings"
	"strconv"
)



func NetStat()(map[string]map[string]uint64, error)  {
	nsFile := "/proc/net/netstat"

	contents, err := ioutil.ReadFile(nsFile)
	if err != nil{
		return nil, fmt.Errorf("read %s failed: %v", nsFile, err)
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	ns := make(map[string]map[string]uint64)


	for{
		line, err := reader.ReadString('\n')
		if err == io.EOF{
			err = nil
			break
		}else if err != nil{
			return nil, fmt.Errorf("read %s buffer failed: %v", nsFile, err)
		}

		idx := strings.Index(line, ":")
		if idx < 0{
			continue
		}

		title := strings.TrimSpace(line[:idx])
		titleMap := ns[title]
		keys := strings.Fields(strings.TrimSpace(line[idx+1:]))
		valLine, err := reader.ReadString('\n')
		if err != nil{
			return nil, fmt.Errorf("process %s buffer for parse %s failed: %v", nsFile, title, err)
		}

		values := strings.Fields(strings.TrimSpace(valLine[idx+1:]))

		for i := 0; i < len(values); i++{
			titleMap[strings.TrimSpace(keys[i])], err = strconv.ParseUint(
				strings.TrimSpace(values[i]), 10, 64)

			if err != nil{
				return ns, fmt.Errorf("process %s buffer for parse %s.%s failed: %v",
					nsFile, title, keys[i], err)
			}

			}


	}

	return ns, nil
}