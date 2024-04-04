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

	for _, group := range ping.Servers {
		if group.Name != params.ServerGroup {
			continue
		}
		for _, server := range group.List {
			statistics[server.Name] = ping.Statistics{
				ServerName: server.Name,
			}
		}
	}

	for i := 0; i < params.RequestCount; i++ {
		fmt.Println("")
		fmt.Printf("Request # %v\n", i+1)

		for _, group := range ping.Servers {
			if group.Name != params.ServerGroup {
				continue
			}
			for _, server := range group.List {
				stat := statistics[server.Name]

				result, err := ping.OpenConnection(server.Host, server.Port, params.Timeout)

				if err == nil && result.Status == "success" {
					stat.ResponseDurations = append(stat.ResponseDurations, int(result.ResponseDuration.Milliseconds()))
					fmt.Println(server.Name, result.ResponseDuration)
				}

				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
						stat.Timeouts++
						fmt.Println(server.Name, "timeout")
					} else {
						stat.Errors++
						fmt.Println(server.Name, err)
					}
				}

				statistics[server.Name] = stat
			}
		}
	}

	ping.PrintResults(statistics)

	fmt.Println("")
	fmt.Print("Press Enter for exit...")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
