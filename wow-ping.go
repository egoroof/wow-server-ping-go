package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/egoroof/wow-server-ping/pkg/ping"
)

var LISTEN_PORT = flag.Int("port", 0, "listen port for Prometheus metrics")
var PING_INTERVAL = flag.Duration("interval", time.Millisecond*500, "sleep time between requests")
var PING_TIMEOUT = flag.Duration("timeout", time.Second, "ping timeout")
var STATS_INTERVAL = flag.Duration("stats-interval", time.Second*30, "how often stats should be printed to console")
var STATS_COUNT = flag.Int("stats", 0, "how many stats to display before exit")
var SERVER_CONFIG = flag.String("servers", "x1", "server config")

var promRespTime = ping.PrometheusMetric{
	Name:       "wow_server_response_time_ms",
	Help:       "WoW server response time in ms",
	Type:       "gauge",
	LabelNames: []string{"server"},
}
var promRespTimeout = ping.PrometheusMetric{
	Name:       "wow_server_timeout_count",
	Help:       "WoW server timeout count",
	Type:       "counter",
	LabelNames: []string{"server"},
}
var promRespErr = ping.PrometheusMetric{
	Name:       "wow_server_error_count",
	Help:       "WoW server error count",
	Type:       "counter",
	LabelNames: []string{"server"},
}

func recordMetrics(servers []ping.Server) {
	responseChan := make(chan ping.ServerResponse)

	for _, server := range servers {
		promRespTimeout.SetValue([]string{server.Name}, 0)
		promRespErr.SetValue([]string{server.Name}, 0)
	}

	statsLogTime := time.Now()
	statistics := make(map[string]ping.Statistics)
	statsCount := 0
	for {
		for _, server := range servers {
			go ping.OpenConnection(server.Name, server.Ip, server.Port, *PING_TIMEOUT, responseChan)
		}

		for range servers {
			response := <-responseChan

			stat, statExist := statistics[response.Name]
			if !statExist {
				stat = ping.Statistics{
					ServerName: response.Name,
				}
			}

			if response.Error == nil {
				promRespTime.SetValue([]string{response.Name}, response.Duration)
				stat.ResponseDurations = append(stat.ResponseDurations, response.Duration)
			} else {
				promRespTime.Delete([]string{response.Name})
				if errors.Is(response.Error, context.DeadlineExceeded) || errors.Is(response.Error, os.ErrDeadlineExceeded) {
					promRespTimeout.AddValue([]string{response.Name}, 1)
					stat.Timeouts++
				} else {
					fmt.Printf("%v %v\n", response.Name, response.Error)
					promRespErr.AddValue([]string{response.Name}, 1)
					stat.Errors++
				}
			}

			statistics[response.Name] = stat
		}

		if time.Now().After(statsLogTime.Add(*STATS_INTERVAL)) {
			fmt.Printf(
				"\n%v to %v\n",
				statsLogTime.Format(time.TimeOnly),
				time.Now().Format(time.TimeOnly),
			)
			ping.PrintResults(statistics)
			statsLogTime = time.Now()
			statistics = make(map[string]ping.Statistics)
			statsCount++

			if *STATS_COUNT == statsCount {
				fmt.Println("Exiting due to stats count flag is set and reached")
				os.Exit(0)
			}
		}

		time.Sleep(*PING_INTERVAL)
	}
}

func main() {
	flag.Parse()
	serversPath := fmt.Sprintf("./servers/%v.json", *SERVER_CONFIG)

	if *STATS_COUNT != 0 {
		fmt.Printf("Stats count is %v\n", *STATS_COUNT)
	}
	fmt.Printf("Ping timeout %v\n", *PING_TIMEOUT)
	fmt.Printf("Ping interval %v\n", *PING_INTERVAL)
	fmt.Printf("Stats interval %v\n", *STATS_INTERVAL)
	fmt.Printf("Server list %v\n", serversPath)

	serversFile, err := os.ReadFile(serversPath)
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		os.Exit(1)
	}

	var servers []ping.Server
	err = json.Unmarshal(serversFile, &servers)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %v servers:\n", len(servers))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, server := range servers {
		fmt.Fprintf(w, "%v\t%v\t%v\n", server.Name, server.Ip, server.Port)
	}
	w.Flush()

	if *LISTEN_PORT == 0 {
		fmt.Println("Listen port is not set. Prometheus metrics disabled")
		recordMetrics(servers)
	} else {
		metrics := []*ping.PrometheusMetric{
			&promRespErr,
			&promRespTime,
			&promRespTimeout,
		}
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			var resp strings.Builder
			for _, metric := range metrics {
				resp.WriteString(metric.GetString())
			}
			w.Write([]byte(resp.String()))
		})

		go recordMetrics(servers)
		fmt.Printf("Listening port %v\n", *LISTEN_PORT)
		err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", *LISTEN_PORT), nil)
		fmt.Println(err)
	}
}
