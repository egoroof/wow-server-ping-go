package ping

import (
	"fmt"
	"os"
	"slices"
	"text/tabwriter"
)

type Statistics struct {
	ServerName        string
	ResponseDurations []int
	Errors            int
	Timeouts          int

	Avg    int
	Jitter int
}

func PrintResults(statistics map[string]Statistics) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	var serverTable []Statistics
	for _, stats := range statistics {
		stats.Avg = Avg(stats.ResponseDurations)
		stats.Jitter = Jitter(stats.ResponseDurations)

		serverTable = append(serverTable, stats)
	}
	slices.SortFunc(serverTable, func(a, b Statistics) int {
		if a.Errors-b.Errors != 0 {
			return a.Errors - b.Errors
		}
		if a.Timeouts-b.Timeouts != 0 {
			return a.Timeouts - b.Timeouts
		}
		return a.Avg - b.Avg
	})

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Response time, ms")
	for _, stats := range serverTable {
		timeoutStr := ""
		if stats.Timeouts > 0 {
			timeoutStr = fmt.Sprintf("; %v timeouts", stats.Timeouts)
		}
		errorStr := ""
		if stats.Errors > 0 {
			errorStr = fmt.Sprintf("; %v errors", stats.Errors)
		}
		statsStr := "unavailable"
		if stats.Avg > 0 {
			statsStr = fmt.Sprintf("%v ± %v", stats.Avg, stats.Jitter)
		}
		fmt.Fprintf(w, "%v\t%v%v%v\n", stats.ServerName, statsStr, timeoutStr, errorStr)
	}
	w.Flush()
}
