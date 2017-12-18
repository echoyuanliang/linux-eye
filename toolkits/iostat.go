package toolkits

import (
	"strconv"
	"io/ioutil"
	"bufio"
	"io"
	"strings"
	"bytes"
	"fmt"
)


type DiskStats struct {
	Major             int			`json:"major"`
	Minor             int			`json:"minor"`
	Device            string		`json:"device"`
	ReadRequests      uint64		`json:"read_requests"`
	ReadMerged        uint64		`json:"read_merged"`
	ReadSectors       uint64		`json:"read_sectors"`
	MsecRead          uint64		`json:"msec_read"`
	WriteRequests     uint64		`json:"write_requests"`
	WriteMerged       uint64		`json:"write_merged"`
	WriteSectors      uint64		`json:"write_sectors"`
	MsecWrite         uint64		`json:"msec_write"`
	IosInProgress     uint64		`json:"ios_in_progress"`
	MsecTotal         uint64		`json:"msec_total"`
	MsecWeightedTotal uint64		`json:"msec_weighted_total"`
}


func ListDiskStats() ([]*DiskStats, error) {
	diskStatsFile := "/proc/diskstats"
	contents, err := ioutil.ReadFile(diskStatsFile)
	if err != nil {
		return nil, fmt.Errorf("read %s failed: %v", diskStatsFile, err)
	}

	ret := make([]*DiskStats, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		fields := strings.Fields(line)
		if fields[3] == "0" {
			continue
		}

		size := len(fields)
		if size != 14 {
			continue
		}

		for idx, field := range fields{
			fields[idx] = strings.TrimSpace(field)
		}

		item := &DiskStats{}
		if item.Major, err = strconv.Atoi(fields[0]); err != nil {
			return nil, err
		}

		if item.Minor, err = strconv.Atoi(fields[1]); err != nil {
			return nil, err
		}

		item.Device = fields[2]

		if item.ReadRequests, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
			return nil, err
		}

		if item.ReadMerged, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
			return nil, err
		}

		if item.ReadSectors, err = strconv.ParseUint(fields[5], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecRead, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
			return nil, err
		}

		if item.WriteRequests, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
			return nil, err
		}

		if item.WriteMerged, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
			return nil, err
		}

		if item.WriteSectors, err = strconv.ParseUint(fields[9], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecWrite, err = strconv.ParseUint(fields[10], 10, 64); err != nil {
			return nil, err
		}

		if item.IosInProgress, err = strconv.ParseUint(fields[11], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecTotal, err = strconv.ParseUint(fields[12], 10, 64); err != nil {
			return nil, err
		}

		if item.MsecWeightedTotal, err = strconv.ParseUint(fields[13], 10, 64); err != nil {
			return nil, err
		}

		ret = append(ret, item)
	}
	return ret, nil
}
