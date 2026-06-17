package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/egoroof/wow-server-ping/pkg/ping"
)

var REQUEST_COUNT = flag.Int("n", 4, "request count")
var PING_INTERVAL = flag.Duration("i", time.Millisecond*500, "sleep time between requests")
var PING_TIMEOUT = flag.Duration("t", time.Second, "ping timeout")
var SERVER_CONFIG = flag.String("s", "x1", "server config")

func main() {
	flag.Parse()
	serversPath := fmt.Sprintf("./servers/%v.json", *SERVER_CONFIG)

	fmt.Printf("Request count %v\n", *REQUEST_COUNT)
	fmt.Printf("Timeout %v\n", *PING_TIMEOUT)
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
	for _, server := range servers {
		fmt.Printf("%v:%v - %v\n", server.Ip, server.Port, server.Name)
	}

	statistics := make(map[string]ping.Statistics)
	responseChan := make(chan ping.ServerResponse)

	for i := 0; i < *REQUEST_COUNT; i++ {
		fmt.Println("")
		fmt.Printf("Request # %v\n", i+1)

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
				stat.ResponseDurations = append(stat.ResponseDurations, response.Duration)
				fmt.Printf("%v %vms\n", response.Name, response.Duration)
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

		time.Sleep(*PING_INTERVAL)
	}

	ping.PrintResults(statistics)
}
