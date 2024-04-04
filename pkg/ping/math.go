package ping

import (
	"math"
)

func Avg(values []int) int {
	if len(values) == 0 {
		return 0
	}

	sum := 0
	for _, elem := range values {
		sum += elem
	}

	return sum / len(values)
}

func Jitter(values []int) int {
	if len(values) == 0 {
		return 0
	}

	min := math.MaxInt
	max := 0
	for _, value := range values {
		if min > value {
			min = value
		}
		if max < value {
			max = value
		}
	}

	return (max - min) / 2
}
