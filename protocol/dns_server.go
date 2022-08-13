package protocol

import (
	"DnsLog/config"
	"DnsLog/store"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/dns/dnsmessage"
	"net"
	"strings"
	"time"
)

type DnsInfo struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"sub_domain"`
	Address   string `json:"host"`
	Time      int64  `json:"time"`
}

// ListingDnsServer 监听dns端口
func ListingDnsServer(options *config.Options) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: options.DNSPort})
	if err != nil {
		logrus.Fatal(err.Error())
	}
	defer conn.Close()
	logrus.Infof("DNS Listing Start On %d", options.DNSPort)
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Println(err)
			continue
		}
		go serverDNS(options, addr, conn, msg)
	}
}

func serverDNS(options *config.Options, addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	if len(msg.Questions) < 1 {
		return
	}
	question := msg.Questions[0]
	var (
		queryNameStr = question.Name.String()
		queryType    = question.Type
		queryName, _ = dnsmessage.NewName(queryNameStr)
		resource     dnsmessage.Resource
	)
	//域名过滤
	queryNameStr = strings.TrimRight(queryNameStr, ".")
	logrus.Infof("DNS LOG %s", queryNameStr)
	subNameLen := len(queryNameStr) - len(options.Domain)
	if subNameLen < 0 {
		logrus.Warnf("GET ERROR DNS %s", queryNameStr)
	} else {
		subDomain := queryNameStr[:subNameLen]
		subDomain = strings.TrimRight(subDomain, ".")
		ds := strings.Split(subDomain, ".")
		key := ds[len(ds)-1]
		if strings.Contains(queryNameStr, options.Domain) {
			logrus.Infof("NSLOOK:%s", queryNameStr)
			item, _ := json.Marshal(DnsInfo{
				Domain:    queryNameStr,
				Subdomain: subDomain,
				Address:   addr.IP.String(),
				Time:      time.Now().Unix(),
			})
			logrus.Infof("DNS SET %s %s", key, item)
			err := store.Store.SetItem(key, string(item))
			if err != nil {
				logrus.Warnf(err.Error())
				return
			}
		}
	}
	switch queryType {
	case dnsmessage.TypeA:
		resource = NewAResource(queryName, [4]byte{127, 0, 0, 1})
	default:
		resource = NewAResource(queryName, [4]byte{127, 0, 0, 1})
	}
	// send response
	msg.Response = true
	msg.Answers = append(msg.Answers, resource)
	Response(addr, conn, msg)
}

// Response return
func Response(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	packed, err := msg.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := conn.WriteToUDP(packed, addr); err != nil {
		fmt.Println(err)
	}
}

func NewAResource(query dnsmessage.Name, a [4]byte) dnsmessage.Resource {
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  query,
			Class: dnsmessage.ClassINET,
			TTL:   0,
		},
		Body: &dnsmessage.AResource{
			A: a,
		},
	}
}
