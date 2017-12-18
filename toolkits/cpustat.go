package toolkits

import (
	"strconv"
	"io/ioutil"
	"bufio"
	"io"
	"runtime"
	"bytes"
	"strings"
	"errors"
	"fmt"
)


type CpuUsage struct {
	User    uint64	`json:"user"`
	Nice    uint64	`json:"nice"`
	System  uint64	`json:"system"`
	Idle    uint64	`json:"idle"`
	IoWait  uint64	`json:"io_wait"`
	Irq     uint64	`json:"irq"`
	SoftIrq uint64	`json:"soft_irq"`
	Steal   uint64	`json:"steal"`
	Guest   uint64	`json:"guest"`
	Total   uint64	`json:"total"`
}


type ProcStat struct {
	Cpu          *CpuUsage	`json:"cpu"`
	CpuList         []*CpuUsage	`json:"cpu_list"`
	Ctxt         uint64	`json:"ctxt"`
	Processes    uint64	`json:"processes"`
	ProcessRunning uint64	`json:"process_running"`
	ProcessBlocked uint64		`json:"process_blocked"`
}


func CurrentProcStat() (*ProcStat, error) {
	f := "/proc/stat"
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	ps := &ProcStat{CpuList: make([]*CpuUsage, runtime.NumCPU())}
	reader := bufio.NewReader(bytes.NewBuffer(bs))

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			errMsg := fmt.Sprintf("read /proc/stat buffer failed: %v", err)
			log.Error(errMsg)
			return ps, errors.New(errMsg)
		}
		parseLine(line, ps)
	}

	return ps, nil
}

func parseLine(line string, ps *ProcStat) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return
	}

	fieldName := fields[0]
	if fieldName == "cpu" {
		ps.Cpu = parseCpuFields(fields)
		return
	}

	if strings.HasPrefix(fieldName, "cpu") {
		idx, err := strconv.Atoi(fieldName[3:])
		if err != nil || idx >= len(ps.CpuList) {
			return
		}

		ps.CpuList[idx] = parseCpuFields(fields)
		return
	}

	if fieldName == "ctxt" {
		ps.Ctxt, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "processes" {
		ps.Processes, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "procs_running" {
		ps.ProcessRunning, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "procs_blocked" {
		ps.ProcessBlocked, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
}

func parseCpuFields(fields []string) *CpuUsage {
	cu := new(CpuUsage)
	sz := len(fields)
	for i := 1; i < sz; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		cu.Total += val
		switch i {
		case 1:
			cu.User = val
		case 2:
			cu.Nice = val
		case 3:
			cu.System = val
		case 4:
			cu.Idle = val
		case 5:
			cu.IoWait = val
		case 6:
			cu.Irq = val
		case 7:
			cu.SoftIrq = val
		case 8:
			cu.Steal = val
		case 9:
			cu.Guest = val
		}
	}
	return cu
}

