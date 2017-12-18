package toolkits

import (
	"io"
	"fmt"
	"bytes"
	"bufio"
	"io/ioutil"
	"strconv"
	"strings"
)


const (
	BitsPerByte = 8
	MillionBit   = 1000000
)

type NetIf struct {
	Iface          string		`json:"iface"`
	InBytes        int64		`json:"in_bytes"`
	InPackages     int64		`json:"in_packages"`
	InErrors       int64		`json:"in_errors"`
	InDropped      int64		`json:"in_dropped"`
	InFifoErrs     int64		`json:"in_fifo_errs"`
	InFrameErrs    int64		`json:"in_frame_errs"`
	InCompressed   int64		`json:"in_compressed"`
	InMulticast    int64		`json:"in_multicast"`
	OutBytes       int64		`json:"out_bytes"`
	OutPackages    int64		`json:"out_packages"`
	OutErrors      int64		`json:"out_errors"`
	OutDropped     int64		`json:"out_dropped"`
	OutFifoErrs    int64		`json:"out_fifo_errs"`
	OutCollisions  int64		`json:"out_collisions"`
	OutCarrierErrs int64		`json:"out_carrier_errs"`
	OutCompressed  int64		`json:"out_compressed"`
	TotalBytes     int64		`json:"total_bytes"`
	TotalPackages  int64		`json:"total_packages"`
	TotalErrors    int64		`json:"total_errors"`
	TotalDropped   int64		`json:"total_dropped"`
	SpeedBits      int64		`json:"speed_bits"`
	InPercent      float64	`json:"in_percent"`
	OutPercent     float64	`json:"out_percent"`
}

func NetIfs() ([]*NetIf, error) {
	netDevFile := "/proc/net/dev"
	contents, err := ioutil.ReadFile(netDevFile)
	if err != nil {
		return nil, fmt.Errorf("read %s failed; %v", netDevFile, err)
	}

	ret := make([]*NetIf, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, fmt.Errorf("read %s buffer failed; %v", netDevFile, err)
		}

		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}

		netIf := NetIf{}

		eth := strings.TrimSpace(line[0:idx])

		netIf.Iface = eth

		fields := strings.Fields(line[idx+1:])

		if len(fields) != 16 {
			continue
		}

		for idx, field := range fields{
			fields[idx] = strings.TrimSpace(field)
		}

		netIf.InBytes, _ = strconv.ParseInt(fields[0], 10, 64)
		netIf.InPackages, _ = strconv.ParseInt(fields[1], 10, 64)
		netIf.InErrors, _ = strconv.ParseInt(fields[2], 10, 64)
		netIf.InDropped, _ = strconv.ParseInt(fields[3], 10, 64)
		netIf.InFifoErrs, _ = strconv.ParseInt(fields[4], 10, 64)
		netIf.InFrameErrs, _ = strconv.ParseInt(fields[5], 10, 64)
		netIf.InCompressed, _ = strconv.ParseInt(fields[6], 10, 64)
		netIf.InMulticast, _ = strconv.ParseInt(fields[7], 10, 64)

		netIf.OutBytes, _ = strconv.ParseInt(fields[8], 10, 64)
		netIf.OutPackages, _ = strconv.ParseInt(fields[9], 10, 64)
		netIf.OutErrors, _ = strconv.ParseInt(fields[10], 10, 64)
		netIf.OutDropped, _ = strconv.ParseInt(fields[11], 10, 64)
		netIf.OutFifoErrs, _ = strconv.ParseInt(fields[12], 10, 64)
		netIf.OutCollisions, _ = strconv.ParseInt(fields[13], 10, 64)
		netIf.OutCarrierErrs, _ = strconv.ParseInt(fields[14], 10, 64)
		netIf.OutCompressed, _ = strconv.ParseInt(fields[15], 10, 64)

		netIf.TotalBytes = netIf.InBytes + netIf.OutBytes
		netIf.TotalPackages = netIf.InPackages + netIf.OutPackages
		netIf.TotalErrors = netIf.InErrors + netIf.OutErrors
		netIf.TotalDropped = netIf.InDropped + netIf.OutDropped

		speedFile := fmt.Sprintf("/sys/class/net/%s/speed", netIf.Iface)
		if content, err := ioutil.ReadFile(speedFile); err == nil {
			var speed int64
			speed, err = strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
			if err != nil || speed == 0 {
				netIf.SpeedBits = int64(0)
				netIf.InPercent = float64(0)
				netIf.OutPercent = float64(0)
			} else {
				netIf.SpeedBits = speed * MillionBit
				netIf.InPercent = float64(netIf.InBytes* BitsPerByte) * 100.0 / float64(netIf.SpeedBits)
				netIf.OutPercent = float64(netIf.OutBytes* BitsPerByte) * 100.0 / float64(netIf.SpeedBits)
			}
		} else {
			netIf.SpeedBits = int64(0)
			netIf.InPercent = float64(0)
			netIf.OutPercent = float64(0)
		}

		ret = append(ret, &netIf)
	}

	return ret, nil
}
