package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/egoroof/wow-server-ping-go/pkg/ping"
)

func main() {
	params := ping.ParseArguments(os.Args[1:])

	fmt.Printf("Requests limit %v\n", params.RequestCount)
	fmt.Printf("Timeout %v ms\n", params.Timeout)
	fmt.Printf("Server group '%v'\n", params.ServerGroup)
	statistics := make(map[string]ping.Statistics)
	responseChan := make(chan ping.ServerResponse)

	for i := 0; i < params.RequestCount; i++ {
		fmt.Println("")
		fmt.Printf("Request # %v\n", i+1)
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
