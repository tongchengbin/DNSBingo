package api

import (
	"DnsLog/config"
	"embed"
	"fmt"
	"log"
	"net/http"
)

//go:embed template
var template embed.FS

func ListingHttpManagementServer() {
	mux := http.NewServeMux()
	//mux.Handle("/template/", http.FileServer(http.FS(template)))
	//mux.HandleFunc("/", index)
	mux.HandleFunc("/api/getDomain", getDomain)
	mux.HandleFunc("/api/getRecords", getRecords)
	log.Println("API Listing Start...")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
