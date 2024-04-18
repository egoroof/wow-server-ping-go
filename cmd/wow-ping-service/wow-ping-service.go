package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/egoroof/wow-server-ping-go/pkg/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = 8090
const sleepBetweenRequestsMs = 500

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
	params := ping.ParseArguments(os.Args[1:])
	responseChan := make(chan ping.ServerResponse)

	for _, group := range ping.Servers {
		if group.Name != params.ServerGroup {
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
			if group.Name != params.ServerGroup {
				continue
			}

			for _, server := range group.List {
				connectionCount++
				go ping.OpenConnection(server.Name, server.Host, server.Port, params.Timeout, responseChan)
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

		time.Sleep(time.Millisecond * sleepBetweenRequestsMs)
	}
}

func main() {
	promReg := prometheus.NewRegistry()
	promReg.MustRegister(promRespTime)
	promReg.MustRegister(promRespTimeout)
	promReg.MustRegister(promRespErr)

	handler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	go recordMetrics()

	log.Printf("Listening port %v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	log.Fatal(err)
}
