package main

import (
	"context"
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
var SERVER_GROUP = flag.String("s", "x1", "server group")

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

func recordMetrics() {
	responseChan := make(chan ping.ServerResponse)

	for _, group := range ping.Servers {
		if group.Name != *SERVER_GROUP {
			continue
		}

		for _, server := range group.List {
			promRespTimeout.WithLabelValues(server.Name).Add(0)
			promRespErr.WithLabelValues(server.Name).Add(0)
		}
	}

	for {
		connectionCount := 0

		for _, group := range ping.Servers {
			if group.Name != *SERVER_GROUP {
				continue
			}

			for _, server := range group.List {
				connectionCount++
				go ping.OpenConnection(server.Name, server.Host, server.Port, *TIMEOUT, responseChan)
			}
		}

		for i := 0; i < connectionCount; i++ {
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
	log.Printf("Timeout %v ms\n", *TIMEOUT)
	log.Printf("Server group '%v'\n", *SERVER_GROUP)

	promReg := prometheus.NewRegistry()
	promReg.MustRegister(promRespTime)
	promReg.MustRegister(promRespTimeout)
	promReg.MustRegister(promRespErr)

	handler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	go recordMetrics()

	log.Printf("Listening port %v\n", *PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *PORT), nil)
	log.Fatal(err)
}
