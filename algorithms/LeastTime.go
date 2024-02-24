package algorithms

import (
	constants "github.com/lokesh-katari/GoScale/constants"
	// "time"
)

func LeastTime(backends []*constants.Backend) int {
	var min int
	min = int(backends[0].ResponseTime)
	index := 0
	for i := 0; i < len(backends); i++ {
		if int(backends[i].ResponseTime) < min {
			min = int(backends[i].ResponseTime)
			index = i
		}
	}
	return index

}
