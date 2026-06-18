package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"text/tabwriter"
	"time"

	"github.com/egoroof/wow-server-ping/pkg/ping"
)

var REQUEST_COUNT = flag.Int("requests", 4, "request count")
var PING_INTERVAL = flag.Duration("interval", time.Millisecond*500, "sleep time between requests")
var PING_TIMEOUT = flag.Duration("timeout", time.Second, "ping timeout")
var SERVER_CONFIG = flag.String("servers", "x1", "server config")

func main() {
	flag.Parse()
	serversPath := fmt.Sprintf("./servers/%v.json", *SERVER_CONFIG)

	fmt.Printf("Request count %v\n", *REQUEST_COUNT)
	fmt.Printf("Timeout %v\n", *PING_TIMEOUT)
	fmt.Printf("Interval %v\n", *PING_INTERVAL)
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
	fmt.Fprintf(w, "Name\tIP\tPort\n")
	for _, server := range servers {
		fmt.Fprintf(w, "%v\t%v\t%v\n", server.Name, server.Ip, server.Port)
	}
	w.Flush()

	statistics := make(map[string]ping.Statistics)
	responseChan := make(chan ping.ServerResponse)

	for i := 0; i < *REQUEST_COUNT; i++ {
		fmt.Println("")
		fmt.Printf("Request # %v\n", i+1)

		for _, server := range servers {
			go ping.OpenConnection(server.Name, server.Ip, server.Port, *PING_TIMEOUT, responseChan)
		}

		var responses []ping.ServerResponse
		for range servers {
			response := <-responseChan
			responses = append(responses, response)

			stat, statExist := statistics[response.Name]
			if !statExist {
				stat = ping.Statistics{
					ServerName: response.Name,
				}
			}

			if response.Error == nil {
				stat.ResponseDurations = append(stat.ResponseDurations, response.Duration)
			} else {
				if errors.Is(response.Error, context.DeadlineExceeded) || errors.Is(response.Error, os.ErrDeadlineExceeded) {
					stat.Timeouts++
					fmt.Println(response.Name, "timeout")
				} else {
					stat.Errors++
					fmt.Println(response.Name, response.Error)
				}
			}

			statistics[response.Name] = stat
		}
		slices.SortFunc(responses, func(a, b ping.ServerResponse) int {
			return a.Duration - b.Duration
		})

		for _, response := range responses {
			if response.Duration == 0 {
				continue
			}
			fmt.Fprintf(w, "%v\t%vms\n", response.Name, response.Duration)
		}
		w.Flush()

		time.Sleep(*PING_INTERVAL)
	}

	ping.PrintResults(statistics)
}
