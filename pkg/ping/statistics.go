package ping

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

type Statistics struct {
	ServerName        string
	ServerGroup       string
	ResponseDurations []int
	Errors            int
	Timeouts          int

	Avg    int
	Jitter int
}

func PrintResults(statistics map[string]Statistics, groupsOrder string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	groups := strings.Split(groupsOrder, ",")

	serverTableGroups := make(map[string][]Statistics)
	for _, stats := range statistics {
		stats.Avg = Avg(stats.ResponseDurations)
		stats.Jitter = Jitter(stats.ResponseDurations)

		serverTableGroups[stats.ServerGroup] = append(serverTableGroups[stats.ServerGroup], stats)
	}
	for group := range serverTableGroups {
		slices.SortFunc(serverTableGroups[group], func(a, b Statistics) int {
			if a.Errors-b.Errors != 0 {
				return a.Errors - b.Errors
			}
			if a.Timeouts-b.Timeouts != 0 {
				return a.Timeouts - b.Timeouts
			}
			return a.Avg - b.Avg
		})
	}

	for _, group := range groups {
		for _, stats := range serverTableGroups[group] {
			timeoutStr := ""
			if stats.Timeouts > 0 {
				timeoutStr = fmt.Sprintf("%v timeouts", stats.Timeouts)
			}
			errorStr := ""
			if stats.Errors > 0 {
				errorStr = fmt.Sprintf("%v errors", stats.Errors)
			}
			statsStr := "unavailable"
			if stats.Avg > 0 {
				statsStr = fmt.Sprintf("%v\t± %v", stats.Avg, stats.Jitter)
			}
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", stats.ServerName, statsStr, timeoutStr, errorStr)
		}
		w.Flush()
		if len(serverTableGroups) > 1 {
			fmt.Println("")
		}
	}
}
