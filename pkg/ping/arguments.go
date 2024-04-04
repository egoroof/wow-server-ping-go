package ping

import (
	"strconv"
	"strings"
)

const REQUEST_COUNT = 4
const TIMEOUT = 1000
const SERVER_GROUP = "x1"

type Params struct {
	RequestCount int
	Timeout      int
	ServerGroup  string
}

func parseIntOrDefault(value string, def int) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}

	return i
}

// hard coded
func ParseArguments(args []string) Params {
	params := Params{
		RequestCount: REQUEST_COUNT,
		Timeout:      TIMEOUT,
		ServerGroup:  SERVER_GROUP,
	}
	for _, arg := range args {
		// -n=1 -t=100 -s=Fun
		parts := strings.Split(arg, "=")
		if len(parts) != 2 || len(parts[0]) < 2 {
			continue
		}

		param := parts[0][1:] // -n -> n
		value := parts[1]

		if param == "n" {
			params.RequestCount = parseIntOrDefault(value, REQUEST_COUNT)
		}

		if param == "t" {
			params.Timeout = parseIntOrDefault(value, TIMEOUT)
		}

		if param == "s" {
			for _, group := range Servers {
				if group.Name == value {
					params.ServerGroup = value
					break
				}
			}

		}
	}

	return params
}
