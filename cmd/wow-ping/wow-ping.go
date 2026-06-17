package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/egoroof/wow-server-ping-go/pkg/ping"
)

var REQUEST_COUNT = flag.Int("n", 6, "request count")
var TIMEOUT = flag.Int("t", 1000, "timeout")
var SERVER_CONFIG = flag.String("s", "x1", "server config")

func main() {
	flag.Parse()
	serversPath := fmt.Sprintf("./servers/%v.json", *SERVER_CONFIG)

	fmt.Printf("Requests limit %v\n", *REQUEST_COUNT)
	fmt.Printf("Timeout %v ms\n", *TIMEOUT)
	fmt.Printf("Servers list %v\n", serversPath)

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

	statistics := make(map[string]ping.Statistics)
	responseChan := make(chan ping.ServerResponse)

	for i := 0; i < *REQUEST_COUNT; i++ {
		fmt.Println("")
		fmt.Printf("Request # %v\n", i+1)

		for _, server := range servers {
			go ping.OpenConnection(server.Name, server.Host, server.Port, *TIMEOUT, responseChan)
		}

		for i := 0; i < len(servers); i++ {
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
	}

	ping.PrintResults(statistics)

	fmt.Println("")
	fmt.Print("Press Enter for exit...")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
