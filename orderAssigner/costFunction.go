package orderAssigner

import (
	"../def"
	"../elevatorLogic"
)

const (
	DOOR_OPEN_DURATION = 3000 // time in milliseconds
	TRAVEL_DURATION    = 2500 // time in milliseconds
)

// Runs a simulation of an elevator with its current orders and returns the time until it is finished
func timeToIdle(simulationElevator def.States) int {

	duration := 0

	switch simulationElevator.Behaviour {

	case def.IDLE:
		simulationElevator.Direction = elevatorLogic.ChooseDirection(simulationElevator)
		if simulationElevator.Direction == def.STOP {
			return duration
		}
	case def.MOVING:
		duration += TRAVEL_DURATION / 2
		simulationElevator.Floor += int(simulationElevator.Direction)
	case def.DOOR_OPEN:
		duration -= DOOR_OPEN_DURATION / 2

	}

	for {
		if elevatorLogic.ShouldStop(simulationElevator) {
			simulationElevator = elevatorLogic.ClearOrderAtCurrentFloor(simulationElevator)
			duration += DOOR_OPEN_DURATION
			simulationElevator.Direction = elevatorLogic.ChooseDirection(simulationElevator)
			// Terminate loop when next direction is "STOP", indicating
			// that there is nowhere to go and that we are idle.
			if simulationElevator.Direction == def.STOP {
				return duration
			}
		}

		// Travel to next floor and add travel time to total duration
		simulationElevator.Floor += int(simulationElevator.Direction)
		duration += TRAVEL_DURATION
	}

}
