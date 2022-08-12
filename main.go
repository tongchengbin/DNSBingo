package main

import (
	"DnsLog/api"
	"DnsLog/config"
	"DnsLog/protocol"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"os"
)

func initLogging() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	cliOptions := &config.CLIServerOptions{LdapPort: 1389}
	flag.StringVar(&cliOptions.Domain, "domain", "", "域名地址")
	flag.IntVar(&cliOptions.ManagePort, "port", 8000, "api管理端口")
	flag.StringVar(&cliOptions.ListenIP, "l", "0.0.0.0", "api 监听IP")
	flag.StringVar(&cliOptions.IPAddress, "host", "127.0.0.1", "主机IP")
	flag.IntVar(&cliOptions.DNSPort, "dns-port", 53, "主机IP")
	flag.Parse()
	if cliOptions.Domain == "" {
		flag.Usage()
		os.Exit(1)
	}
	initLogging()
	options := cliOptions.AsServerOptions()
	config.OptionsConfig = *options
	// DNS
	go protocol.ListingDnsServer(options)

	// LDAP
	serverOptions := cliOptions.AsServerOptions()
	ldapServer, err := protocol.NewLDAPServer(serverOptions, false)
	if err != nil {
		fmt.Printf("Could not create LDAP server: %s", err)
		os.Exit(1)
	}
	var tlsConfig *tls.Config
	go ldapServer.ListenAndServe(tlsConfig)
	defer ldapServer.Close()
	api.ListingHttpManagementServer(options)
}
