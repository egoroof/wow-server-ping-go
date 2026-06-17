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

	"github.com/egoroof/wow-server-ping/pkg/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var LISTEN_PORT = flag.Int("port", 8090, "listen port")
var PING_INTERVAL = flag.Duration("interval", time.Millisecond*500, "sleep time between requests")
var PING_TIMEOUT = flag.Duration("timeout", time.Second, "ping timeout")
var SERVER_CONFIG = flag.String("servers", "x1", "server config")

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
			go ping.OpenConnection(server.Name, server.Ip, server.Port, *PING_TIMEOUT, responseChan)
		}

		for range servers {
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

		time.Sleep(*PING_INTERVAL)
	}
}

func main() {
	flag.Parse()
	serversPath := fmt.Sprintf("./servers/%v.json", *SERVER_CONFIG)

	log.Printf("Timeout %v\n", *PING_TIMEOUT)
	log.Printf("Server list %v\n", serversPath)

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

	log.Printf("Loaded %v servers:\n", len(servers))
	for _, server := range servers {
		log.Printf("%v:%v - %v\n", server.Ip, server.Port, server.Name)
	}

	promReg := prometheus.NewRegistry()
	promReg.MustRegister(promRespTime)
	promReg.MustRegister(promRespTimeout)
	promReg.MustRegister(promRespErr)

	handler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	go recordMetrics(servers)

	log.Printf("Listening port %v\n", *LISTEN_PORT)
	err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", *LISTEN_PORT), nil)
	log.Fatal(err)
}
