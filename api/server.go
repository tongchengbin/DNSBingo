package api

import (
	"DnsLog/config"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

}
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
		logrus.Fatal(err)
	}
}
