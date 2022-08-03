package api

import (
	"DnsLog/config"
	"DnsLog/store"
	"DnsLog/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RespData struct {
	HTTPStatusCode int
	Msg            string
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/template", http.StatusMovedPermanently)
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("domain")
	data, _ := json.Marshal(store.GetDns(key))
	fmt.Fprintf(w, string(data))
}

func getDomain(w http.ResponseWriter, r *http.Request) {
	// get subdomain random id
	domain := utils.RandomString(6)
	fullDomain := domain + "." + config.Domain
	// 如果需要身份验证 这里绑定身份就可以了  、
	store.CreateKey(domain)
	io.WriteString(w, fullDomain)

}
