package tests

import (
	"DnsLog/store"
	"encoding/json"
	"testing"
)

type DNSData struct {
	Protocol string
}

func TestSave(t *testing.T) {
	a := DNSData{"http"}
	marshal, err := json.Marshal(a)
	if err != nil {
		return
	}
	store.Store.RegisterKey("111")
	store.Store.SetItem("111", string(marshal))
	store.Store.SetItem("111", string(marshal))
	data, _ := store.Store.GetItem("111")
	var dat map[string]interface{}
	var results []interface{}
	for _, i := range data.Data {
		err := json.Unmarshal([]byte(i), &dat)
		if err == nil {
			results = append(results, dat)
		}
	}
	f, _ := json.Marshal(results)
	println(string(f))
}
