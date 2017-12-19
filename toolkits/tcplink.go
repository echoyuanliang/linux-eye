package toolkits

import (
	"io/ioutil"
	"fmt"
	"bufio"
	"bytes"
	"io"
	"strings"
	"strconv"
	"os/user"
)

type User struct{
	Uid 		string	`json:"uid"`
	Name 	string	`json:"name"`
	Username		string	`json:"username"`
}


type Address struct{
	Ip 	string	`json:"ip"`
	Port		uint64	`json:"port"`
}


type TcpLink struct{
	LocalAddress		*Address 	`json:"local_address"`
	RemoteAddress	*Address		`json:"remote_address"`

	Status 		string	`json:"status"`
	TxQueue		uint64	`json:"tx_queue"`
	RxQueue		uint64	`json:"rx_queue"`
	TmWhen		uint64	`json:"tm_when"`
	TmRetry		uint64	`json:"tm_retry"`
	User			*User	`json:"user"`
	Inode		uint64	`json:"inode"`
}

var statusList = []string{"ERROR_STATUS",  "TCP_ESTABLISHED",  "TCP_SYN_SENT",  "TCP_SYN_RECV",
	"TCP_FIN_WAIT1",  "TCP_FIN_WAIT2",  "TCP_TIME_WAIT", "TCP_CLOSE", "TCP_CLOSE_WAIT",  "TCP_LAST_ACK",
	"TCP_LISTEN", "TCP_CLOSING"}

func getIpStr(addr uint64) string{

	addrBytes := make([]byte, 4)

	addrBytes[0] = byte(addr & 0xFF)
	addrBytes[1] = byte((addr >> 8) & 0xFF)
	addrBytes[2] = byte((addr >> 16) & 0xFF)
	addrBytes[3] = byte((addr >> 24) & 0xFF)

	return fmt.Sprintf("%d.%d.%d.%d", addrBytes[3], addrBytes[2], addrBytes[1], addrBytes[0])
}


func parseAddress(address string)(*Address, error){
	fields := strings.Split(address, ":")

	if len(fields) != 2{
		return nil, fmt.Errorf("parse address %s failed, invalid format", address)
	}


	a := &Address{}

	port, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 16, 64)
	if err != nil{
		return a, fmt.Errorf("parse address %s failed, invalid port %s ", address, fields[1])
	}

	a.Port = port

	ip, err := strconv.ParseUint(strings.TrimSpace(fields[0]), 16, 64)

	if err != nil{
		return a, fmt.Errorf("parse address %s failed, invalid ip %s", address, fields[0])
	}
	a.Ip = getIpStr(ip)

	return a, nil
}


func parseStatus(status string)(string, error){
	st, err := strconv.ParseUint(strings.TrimSpace(status), 16, 64)
	if err != nil{
		return "", fmt.Errorf("invalid status %s", status)
	}

	idx := int(st)
	if idx > len(statusList){
		return "", fmt.Errorf("invalid status %s, must less than %d", status, idx)
	}

	return statusList[idx], nil
}


func getUserById(uid string)(*User, error){
	if u, err := user.LookupId(strings.TrimSpace(uid)); err != nil{
		return nil, fmt.Errorf("user.LookupId %s failed %v", uid, err)
	}else{
		return &User{Uid:u.Uid, Name:u.Name, Username:u.Username}, nil
	}
}

func GetTcpLinks(tcpFile string)([]*TcpLink, error){
	contents, err := ioutil.ReadFile(tcpFile)

	if err != nil{
		return nil, fmt.Errorf("read %s failed: %v", tcpFile, err)
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	tcpLinks := make([]*TcpLink, 0)

	header := true
	for{
		line, err := reader.ReadString('\n')
		if err == io.EOF{
			err = nil
			break
		}else if err != nil{
			return nil, fmt.Errorf("read %s buffer failed: %v", tcpFile, err)
		}

		link := &TcpLink{}
		if ! header{

			fields := strings.Fields(line)
			if link.LocalAddress, err = parseAddress(fields[1]); err != nil{
				return tcpLinks, fmt.Errorf("parse LocalAddress %s failed: %v", fields[1], err)
			}

			if link.RemoteAddress, err = parseAddress(fields[2]); err != nil{
				return tcpLinks, fmt.Errorf("parse RemoteAddress %s failed: %v", fields[2], err)
			}

			if link.Status, err = parseStatus(fields[3]); err != nil{
				return tcpLinks, fmt.Errorf("parse st %s failed: %v", fields[2], err)
			}

			queueFields := strings.Split(fields[4], ":")
			if len(queueFields) != 2{
				return tcpLinks, fmt.Errorf("parse tx/rx queue %s failed: %v", fields[4], err)
			}

			if link.TxQueue, err = strconv.ParseUint(strings.TrimSpace(queueFields[0]), 16, 64); err != nil{
				return tcpLinks, fmt.Errorf("parse tx queue %s failed: %v", queueFields[0], err)
			}

			if link.RxQueue, err = strconv.ParseUint(strings.TrimSpace(queueFields[1]), 16, 64); err != nil{
				return tcpLinks, fmt.Errorf("parse rx queue %s failed: %v", queueFields[1], err)
			}

			tmFields := strings.Split(fields[5], ":")
			if len(tmFields) != 2{
				return tcpLinks, fmt.Errorf("parse tr/tm-when %s failed: %v", fields[5], err)
			}

			if link.TmWhen, err = strconv.ParseUint(strings.TrimSpace(tmFields[1]), 16, 64); err != nil{
				return tcpLinks, fmt.Errorf("parse tm-when %s failed: %v", tmFields[1], err)
			}

			if link.TmRetry, err = strconv.ParseUint(strings.TrimSpace(fields[6]), 16, 64); err != nil{
				return tcpLinks, fmt.Errorf("parse retrnsmt %s failed: %v", fields[6], err)
			}

			if link.User, err = getUserById(fields[7]); err != nil{
				return tcpLinks, fmt.Errorf("getUserById %s failed: %v", fields[7], err)
			}

			if link.Inode, err = strconv.ParseUint(strings.TrimSpace(fields[9]), 10, 64); err != nil{
				return tcpLinks, fmt.Errorf("parse inode %s failed: %v", fields[9], err)
			}

			tcpLinks = append(tcpLinks, link)

		}

		header = false
	}

	return tcpLinks, nil

}

func TcpLinks()([]*TcpLink, error){
	return GetTcpLinks("/proc/net/tcp")
}