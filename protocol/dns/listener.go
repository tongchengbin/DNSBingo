package Dns

import (
	"DnsLog/config"
	"DnsLog/store"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var rw sync.RWMutex

// ListingDnsServer 监听dns端口
func ListingDnsServer() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	log.Println("DNS Listing Start...")
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Println(err)
			continue
		}
		go serverDNS(addr, conn, msg)
	}
}

func serverDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
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
	subDomain := queryNameStr[:len(queryNameStr)-len(config.Domain)]
	subDomain = strings.TrimRight(subDomain, ".")
	ds := strings.Split(subDomain, ".")
	key := ds[len(ds)-1]
	if strings.Contains(queryNameStr, config.Domain) {
		log.Println("LOOKUP :", queryNameStr)
		store.SetDns(key, store.DnsInfo{
			Domain:    queryNameStr,
			Subdomain: subDomain,
			Address:   addr.IP.String(),
			Time:      time.Now().Unix(),
		})
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
