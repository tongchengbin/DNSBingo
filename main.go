package main

import (
	"DnsLog/api"
	"DnsLog/config"
	Dns "DnsLog/protocol/dns"
	"DnsLog/store"
	"flag"
	"os"
)

func main() {
	flag.StringVar(&config.Domain, "domain", "", "域名地址")
	flag.IntVar(&config.Port, "port", 8000, "api监听端口")
	flag.StringVar(&config.Host, "host", "0.0.0.0", "api 监听IP")
	flag.Parse()
	if config.Domain == "" {
		flag.Usage()
		os.Exit(1)
	}
	go Dns.ListingDnsServer()
	api.ListingHttpManagementServer()
	go store.ExpireMap()
}
