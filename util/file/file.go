package file

import (
	"io/ioutil"
	"strconv"
	"strings"
	"fmt"
)


func ToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read %s failed: %v", filePath, err)
	}
	return string(b), nil
}

func ToTrimString(filePath string) (string, error) {
	str, err := ToString(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

func ToUint64(filePath string) (uint64, error) {
	content, err := ToTrimString(filePath)
	if err != nil {
		return 0, err
	}

	var ret uint64
	if ret, err = strconv.ParseUint(content, 10, 64); err != nil {
		return 0, fmt.Errorf("convert %s content to uint64 failed: %v", filePath, err)
	}
	return ret, nil
}

func ToInt64(filePath string) (int64, error) {
	content, err := ToTrimString(filePath)
	if err != nil {
		return 0, err
	}

	var ret int64
	if ret, err = strconv.ParseInt(content, 10, 64); err != nil {
		return 0, fmt.Errorf("convert %s content to int64 failed: %v", filePath, err)
	}
	return ret, nil
}

