package algorithms

import (
	"fmt"

	"github.com/lokesh-katari/GoScale/constants"
)

func WeightedRoundRobbin(current int, backends []*constants.Backend) (int, int) {
	// fmt.Println("Inside WeightedRoundRobbin")
	totalWeight := 0
	for _, backend := range backends {
		totalWeight += backend.Weight

	}
	current = (current + 1) % totalWeight
	selected := current
	fmt.Println("Selected", selected, current, totalWeight)
	for i, backend := range backends {
		if selected < backend.Weight {
			fmt.Println("Selected in the loop", selected, backend.Weight, totalWeight)
			return i, current
		}
		selected -= backend.Weight
		// fmt.Println("Selected in the loop", selected, current, totalWeight)
	}
	return 0, 0
}
