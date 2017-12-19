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

type Mem struct {
	Buffers   uint64		`json:"buffers"`
	Cached    uint64		`json:"cached"`
	MemTotal  uint64		`json:"mem_total"`
	MemFree   uint64		`json:"mem_free"`
	SwapTotal uint64		`json:"swap_total"`
	SwapUsed  uint64		`json:"swap_used"`
	SwapFree  uint64		`json:"swap_free"`
}


var WantField = map[string]bool{
	"Buffers:": true,
	"Cached:": true,
	"MemTotal:": true,
	"MemFree:": true,
	"SwapTotal:": true,
	"SwapFree:": true,
}

const Kb = 1024


func MemInfo()(*Mem, error)  {
	memFile := "/proc/meminfo"
	contents, err := ioutil.ReadFile(memFile)
	if err != nil{
		return nil, fmt.Errorf("read %s failed: %v", memFile, err)
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	mem := &Mem{}

	for{
		line, err := reader.ReadString('\n')

		if err == io.EOF{
			err = nil
			break
		}else if err != nil{
			return nil, fmt.Errorf("read %s buffer failed: %v", memFile, err)
		}

		fields := strings.Split(line, " ")
		fieldName := strings.TrimSpace(fields[0])
		if _, exists := WantField[fieldName]; exists && len(fields) == 3{
			val, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
			if err != nil{
				return mem, fmt.Errorf("process %s failed, value of %s is not positive integer",
					memFile, fieldName)
			}

			val *= Kb
			switch fieldName{
			case "Buffers:":
				mem.Buffers = val
			case "Cached:":
				mem.Cached = val
			case "MemTotal:":
				mem.MemTotal = val
			case "MemFree:":
				mem.MemFree = val
			case "SwapTotal:":
				mem.SwapTotal = val
			case "SwapFree:":
				mem.SwapFree = val
			}

		}


	}

	mem.SwapUsed = mem.SwapTotal - mem.SwapFree
	return mem, nil
}