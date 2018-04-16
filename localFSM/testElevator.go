package localFSM

import (
	"fmt"

	"../def"
)

func printStatus(tempState def.States) {
	fmt.Println("LocalID: ", tempState.Id)
	fmt.Println("Behaviour: ", tempState.Behaviour, "\t Floor: ", tempState.Floor, "\t direction: ", tempState.Direction, "\t stuck: ", tempState.Stuck)
	fmt.Println("HallRequests: ", tempState.HallRequests)
	fmt.Println("acceptedOrders: ", tempState.AcceptedOrders, "\n")
}
