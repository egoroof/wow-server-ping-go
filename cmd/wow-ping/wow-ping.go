package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/egoroof/wow-server-ping-go/pkg/ping"
)

var REQUEST_COUNT = flag.Int("n", 6, "request count")
var TIMEOUT = flag.Int("t", 1000, "timeout")
var SERVER_GROUP = flag.String("s", "x1", "server group")

func main() {
	flag.Parse()

	fmt.Printf("Requests limit %v\n", *REQUEST_COUNT)
	fmt.Printf("Timeout %v ms\n", *TIMEOUT)
	fmt.Printf("Server group '%v'\n", *SERVER_GROUP)
	statistics := make(map[string]ping.Statistics)
	responseChan := make(chan ping.ServerResponse)

	for i := 0; i < *REQUEST_COUNT; i++ {
		fmt.Println("")
		fmt.Printf("Request # %v\n", i+1)
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
