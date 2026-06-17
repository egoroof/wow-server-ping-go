package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/egoroof/wow-server-ping-go/pkg/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var PORT = flag.Int("p", 8090, "port")
var SLEEP_BETWEEN_REQUESTS_MS = flag.Int("sleep", 500, "sleep time between requests in ms")
var TIMEOUT = flag.Int("t", 1000, "timeout")
var SERVER_CONFIG = flag.String("s", "x1", "server config")

var promRespTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "wow_server_response_time_ms",
	Help: "WoW server response time in ms",
}, []string{"server"})

var promRespTimeout = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "wow_server_timeout_count",
	Help: "WoW server timeout count",
}, []string{"server"})

var promRespErr = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "wow_server_error_count",
	Help: "WoW server error count",
}, []string{"server"})

func recordMetrics(servers []ping.Server) {
	responseChan := make(chan ping.ServerResponse)

	for _, server := range servers {
		promRespTimeout.WithLabelValues(server.Name).Add(0)
		promRespErr.WithLabelValues(server.Name).Add(0)
	}

	for {
		for _, server := range servers {
			go ping.OpenConnection(server.Name, server.Host, server.Port, *TIMEOUT, responseChan)
		}

		for i := 0; i < len(servers); i++ {
			response := <-responseChan

			if response.Error == nil {
				promRespTime.WithLabelValues(response.Name).Set(float64(response.Duration))
			} else {
				promRespTime.DeleteLabelValues(response.Name)
				if errors.Is(response.Error, context.DeadlineExceeded) || errors.Is(response.Error, os.ErrDeadlineExceeded) {
					promRespTimeout.WithLabelValues(response.Name).Inc()
				} else {
					log.Printf("%v %v\n", response.Name, response.Error)
					promRespErr.WithLabelValues(response.Name).Inc()
				}
			}
		}

		time.Sleep(time.Millisecond * time.Duration(*SLEEP_BETWEEN_REQUESTS_MS))
	}
}

func main() {
	flag.Parse()
	serversPath := fmt.Sprintf("./servers/%v.json", *SERVER_CONFIG)

	log.Printf("Timeout %v ms\n", *TIMEOUT)
	log.Printf("Servers list %v\n", serversPath)

	serversFile, err := os.ReadFile(serversPath)
	if err != nil {
		log.Println("Error when opening file: ", err)
		os.Exit(1)
	}

	var servers []ping.Server
	err = json.Unmarshal(serversFile, &servers)
	if err != nil {
		log.Println("Error during Unmarshal(): ", err)
		os.Exit(1)
	}

	promReg := prometheus.NewRegistry()
	promReg.MustRegister(promRespTime)
	promReg.MustRegister(promRespTimeout)
	promReg.MustRegister(promRespErr)

	handler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	go recordMetrics(servers)

	log.Printf("Listening port %v\n", *PORT)
	err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", *PORT), nil)
	log.Fatal(err)
}
