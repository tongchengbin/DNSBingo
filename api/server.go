package api

import (
	"DnsLog/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func ListingHttpManagementServer(options *config.Options) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", register)
	mux.HandleFunc("/api/records", getRecords)
	mux.HandleFunc("/java/", getJavaClass)
	mux.HandleFunc("/jndi/register", registerJNDIClass)
	logrus.Infof("API Listing Start on :%d", config.OptionsConfig.ManagePort)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", options.ManagePort),
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
