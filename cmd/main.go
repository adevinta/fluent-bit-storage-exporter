package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adevinta/fluent-bit-storage-exporter/pkg/client"
	"github.com/adevinta/fluent-bit-storage-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http"
)

func main() {
	var exporterPort int
	var fluentBitPort int
	var fluentBitHost string

	flag.StringVar(&fluentBitHost, "fluent-bit-host", "127.0.0.1", "")
	flag.IntVar(&fluentBitPort, "fluent-bit-port", 2020, "")
	flag.IntVar(&exporterPort, "exporter-port", 8080, "")
	flag.Parse()

	fluentBitClient := client.FluentBitClient{FBHost: fluentBitHost, FBPort: fluentBitPort, HTTPClient: http.Client{}}
	fluentBitCollector := metrics.NewCollector(fluentBitClient)
	prometheus.MustRegister(fluentBitCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter is listening in port 8080....")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", exporterPort), nil))
}
