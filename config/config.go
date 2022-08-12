package config

type CLIServerOptions struct {
	Domain     string
	ManagePort int
	ListenIP   string
	LdapPort   int
	DNSPort    int
	IPAddress  string
}

// Options contains configuration options for the servers
type Options struct {
	Domain string
	// api 管理监听端口
	ManagePort int

	IPAddress string
	// IPAddress is the IP address of the current server.
	ListenIP string
	LdapPort int
	DNSPort  int
}

func (cliServerOptions *CLIServerOptions) AsServerOptions() *Options {
	return &Options{
		IPAddress:  cliServerOptions.IPAddress,
		Domain:     cliServerOptions.Domain,
		ListenIP:   cliServerOptions.ListenIP,
		LdapPort:   cliServerOptions.LdapPort,
		ManagePort: cliServerOptions.ManagePort,
		DNSPort:    cliServerOptions.DNSPort,
	}
}

var OptionsConfig Options
