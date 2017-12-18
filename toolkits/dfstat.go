package toolkits

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"syscall"
)

type Mount struct {
	Device string `json:"device"`
	Point  string `json:"point"`
	FsType string `json:"fs_type"`
}

type DeviceUsage struct {
	Mount *Mount `json:"mount"`

	BlocksAll         uint64  `json:"blocks_all"`
	BlocksUsed        uint64  `json:"blocks_used"`
	BlocksFree        uint64  `json:"blocks_free"`
	BlocksUsedPercent float64 `json:"blocks_used_percent"`
	BlocksFreePercent float64 `json:"blocks_free_percent"`

	InodesAll         uint64  `json:"inodes_all"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
	InodesFreePercent float64 `json:"inodes_free_percent"`

	SizeAll         uint64  `json:"size_all"`
	SizeUsed        uint64  `json:"size_used"`
	SizeFree        uint64  `json:"size_free"`
	SizeUsedPercent float64 `json:"size_used_percent"`
	SizeFreePercent float64 `json:"size_free_percent"`
}

var DeviceIgnore = map[string]bool{
	"none":  true,
	"nodev": true,
}

var FsTypeIgnore = map[string]bool{
	"cgroup":     true,
	"debugfs":    true,
	"devpts":     true,
	"devtmpfs":   true,
	"rpc_pipefs": true,
	"rootfs":     true,
}

var PointPrefixIgnore = []string{
	"/sys",
	"/net",
	"/misc",
	"/proc",
	"/lib",
}

func ignorePoint(point string) bool {
	for _, prefix := range PointPrefixIgnore {
		if strings.HasPrefix(point, prefix) {
			return true
		}
	}

	return false
}

func ListMountPoint() ([]*Mount, error) {
	contents, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		errMsg := fmt.Sprintf("read /proc/mounts failed, error: %v", err)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	mounts := make([]*Mount, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			errMsg := fmt.Sprintf("read /proc/mounts buffer failed, error: %v", err)
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}

		fields := strings.Fields(string(line))

		device := fields[0]
		point := fields[1]
		fsType := fields[2]

		if _, exist := DeviceIgnore[device]; exist {
			continue
		}

		if _, exist := FsTypeIgnore[fsType]; exist {
			continue
		}

		if strings.HasPrefix(fsType, "fuse") {
			continue
		}

		if ignorePoint(point) {
			continue
		}

		if strings.HasPrefix(device, "/dev") {
			deviceFound := false
			for idx := range mounts {
				if mounts[idx].Device == device {
					deviceFound = true
					if len(point) < len(mounts[idx].Point) {
						mounts[idx].Point = point
					}
					break
				}
			}
			if !deviceFound {
				mounts = append(mounts, &Mount{Device: device, Point: point, FsType: fsType})
			}
		} else {
			mounts = append(mounts, &Mount{Device: device, Point: point, FsType: fsType})
		}
	}
	return mounts, nil

}

func BuildDeviceUsage(mount *Mount) (*DeviceUsage, error) {
	usage := &DeviceUsage{Mount: mount}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(mount.Point, &fs)
	if err != nil {
		errMsg := fmt.Sprintf("call Statfs failed: %v", err)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	//blocks
	used := fs.Blocks - fs.Bfree
	usage.BlocksAll = fs.Blocks
	usage.BlocksUsed = used
	usage.BlocksFree = fs.Bavail
	if fs.Blocks == 0 {
		usage.BlocksUsedPercent = 0
		usage.BlocksFreePercent = 0
	} else {
		usage.BlocksUsedPercent = float64(used) * 100.0 / float64(used+fs.Bavail)
		usage.BlocksFreePercent = 100.0 - usage.BlocksUsedPercent
	}

	// inodes
	usage.InodesAll = fs.Files
	if fs.Ffree == math.MaxUint64 {
		usage.InodesFree = 0
		usage.InodesUsed = 0
	} else {
		usage.InodesFree = fs.Ffree
		usage.InodesUsed = fs.Files - fs.Ffree
	}

	if fs.Files == 0 {
		usage.InodesUsedPercent = 0
		usage.InodesFreePercent = 0
	} else {
		usage.InodesUsedPercent = float64(usage.InodesUsed) * 100.0 / float64(usage.InodesAll)
		usage.InodesFreePercent = 100.0 - usage.InodesUsedPercent
	}

	// size
	usage.SizeAll = usage.BlocksAll * uint64(fs.Bsize)
	usage.SizeUsed = usage.BlocksUsed * uint64(fs.Bsize)
	usage.SizeFree = usage.SizeAll - usage.SizeUsed

	if usage.SizeAll == 0 {
		usage.SizeUsedPercent = 0
		usage.SizeFreePercent = 0
	} else {
		usage.SizeUsedPercent = float64(usage.SizeUsed) * 100.0 / float64(usage.SizeAll)
		usage.SizeFreePercent = float64(usage.SizeFree) * 100.0 / float64(usage.SizeAll)
	}

	return usage, nil
}

func ListDeviceUsage() ([]*DeviceUsage, error) {
	mounts, err := ListMountPoint()
	if err != nil {
		errMsg := fmt.Sprintf("list mount point failed: %v", err)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	usages := make([]*DeviceUsage, len(mounts))
	for idx, mount := range mounts {
		usage, err := BuildDeviceUsage(mount)
		if err != nil {
			log.Errorf("build device usage of %s failed: %v", mount.Device, err)
			continue
		}

		usages[idx] = usage
	}

	return usages, nil
}
