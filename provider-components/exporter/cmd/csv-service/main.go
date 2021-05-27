package main

import (
	"fmt"
	"github.com/anishj0shi/csv-exporter/pkg/api"
	"github.com/anishj0shi/csv-exporter/pkg/client"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	config := &api.EnvConfig{}
	err := envconfig.Process("csv", config)
	if err != nil {
		log.Fatal(err.Error())
	}
	odataClient := client.NewOdataClient(config)

	mux := http.NewServeMux()
	mux.HandleFunc("/getCSV", odataClient.HandleCSVExport)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), mux); err != nil {
		log.Fatalf("Unable to Start http Server.\n err: %v", err)
	}
}
