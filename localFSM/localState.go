package localFSM

import (
	"../def"
	"../elevatorLogic"
)

var localState def.States

// Acknowledge hall order or accept cab order
func acknowledgeOrder(button def.Button) {
	if button.Type == def.CAB {
		localState.AcceptedOrders[button.Floor][def.CAB] = true
	} else {
		localState.HallRequests[button.Floor][button.Type] = 1
	}
}

// Accept an order and execute if the elevator is idle
func acceptOrder(motorDirCh chan<- def.MotorDirection, stuckTimerCh chan<- timerOption, doorResetCh chan<- bool, doorLightCh chan<- bool, order def.Button) {
	localState.AcceptedOrders[order.Floor][order.Type] = true

	if localState.Floor == order.Floor && localState.Behaviour != def.MOVING {
		// Elevator has just gotten unstuck, clear orders accepted while already at floor
		localState = elevatorLogic.ClearOrderAtCurrentFloor(localState)
		doorLightCh <- true
		doorResetCh <- true
		localState.Behaviour = def.DOOR_OPEN
	} else if localState.Behaviour == def.IDLE {
		executeOrder(motorDirCh, stuckTimerCh)
	}
}
