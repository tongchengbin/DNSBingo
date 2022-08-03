package store

import (
	"sync"
)

type DnsInfo struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"sub_domain"`
	Address   string `json:"host"`
	Time      int64  `json:"time"`
}

// DnsMap 记录所有的dns 数据
var DnsMap = make(map[string][]DnsInfo)

var rw sync.RWMutex

func CreateKey(token string) {
	rw.Lock()
	DnsMap[token] = []DnsInfo{}
	rw.Unlock()
}

func SetDns(token string, data DnsInfo) {
	rw.Lock()
	if DnsMap[token] == nil {
		DnsMap[token] = []DnsInfo{data}
	} else {
		DnsMap[token] = append(DnsMap[token], data)
	}
	rw.Unlock()
}

func GetDns(token string) []DnsInfo {
	rw.RLock()
	if DnsMap[token] != nil {
		data := DnsMap[token]
		rw.RUnlock()
		return data
	}
	rw.RUnlock()
	return []DnsInfo{}
}
