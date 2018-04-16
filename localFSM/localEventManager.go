package localFSM

import (
	"time"

	"../def"
	"../elevatorLogic"
)

func EventManager(hallReqCh <-chan def.OrderStatus, acceptOrderCh <-chan def.Button,
	statesToNetworkCh chan<- def.States, statesToBackupCh chan<- def.States, statesToOrderAssignerCh chan<- def.States,
	floorSensorCh <-chan int, buttonPressCh chan def.Button, buttonLightCh chan<- def.ButtonLight,
	motorDirCh chan<- def.MotorDirection, doorLightCh chan<- bool, id string, initialFloor int) {

	/***************Initialization**********************/
	localState.Id = id
	localState.Behaviour = def.IDLE
	localState.Direction = def.STOP
	localState.Floor = initialFloor
	/**************************************************/

	doorTimedOutCh := make(chan bool)
	doorResetCh := make(chan bool)
	go doorTimer(doorTimedOutCh, doorResetCh)

	isStuckCh := make(chan bool)
	stuckTimerCh := make(chan timerOption)
	go stuckTimer(isStuckCh, stuckTimerCh)

	// Ticker to send states at regular interval
	ticker := time.NewTicker(def.STATE_TRANSMIT_INTERVAL)

	for {
		select {
		case floor := <-floorSensorCh:
			newFloorReached(buttonLightCh, motorDirCh, doorLightCh, stuckTimerCh, isStuckCh, floor)
			if localState.Behaviour == def.DOOR_OPEN {
				doorResetCh <- true
			}

		case button := <-buttonPressCh:
			if elevatorLogic.ShouldOpenDoor(localState, button) {
				doorLightCh <- true
				localState.Behaviour = def.DOOR_OPEN
				doorResetCh <- true

			} else {
				newOrderPlaced(buttonLightCh, motorDirCh, button, stuckTimerCh)
			}

		case <-doorTimedOutCh:
			doorTimerTimedOut(motorDirCh, doorLightCh, stuckTimerCh)

		case hallReq := <-hallReqCh:
			// Update the hallrequest status with status from orderAssigner
			localState.HallRequests[hallReq.Floor][hallReq.Type] = hallReq.Status

		case order := <-acceptOrderCh:
			// Accept hall request distributed by orderAssigner
			acceptOrder(motorDirCh, stuckTimerCh, doorResetCh, doorLightCh, order)

		case <-ticker.C:
			statesToNetworkCh <- localState
			statesToBackupCh <- localState
			statesToOrderAssignerCh <- localState
		case <-isStuckCh:
			elevatorStuck()

		}

	}

}

// Treat orders at the floor it arrives at or keep going if there are no orders
func newFloorReached(buttonLightCh chan<- def.ButtonLight, motorDirCh chan<- def.MotorDirection,
	doorLightCh chan<- bool, stuckTimerCh chan<- timerOption, isStuckCh chan<- bool, floor int) {
	if localState.Behaviour != def.MOVING {
		return
	}

	localState.Floor = floor

	if !elevatorLogic.ShouldStop(localState) {
		// Elevator continues to move, so start the stuck-timer
		stuckTimerCh <- RESET
		return
	}
	motorDirCh <- def.STOP

	if localState.Stuck {
		// When the elevator reaches a new floor, it is no longer stuck
		localState.Stuck = false

		// Set to idle if there is no cab order
		if !localState.AcceptedOrders[floor][def.CAB] {
			localState.Behaviour = def.IDLE
			return
		}
	} else {
		// Disable stuck timer when idle
		stuckTimerCh <- STOP
	}

	// Treat cab order
	doorLightCh <- true
	localState.Behaviour = def.DOOR_OPEN
	light := def.ButtonLight{def.CAB, floor, false}
	buttonLightCh <- light

	localState = elevatorLogic.ClearOrderAtCurrentFloor(localState)

}

// Acknowledges the order and executes cab orders if the elevator is Idle
func newOrderPlaced(buttonLightCh chan<- def.ButtonLight, motorDirCh chan<- def.MotorDirection,
	button def.Button, stuckTimerCh chan<- timerOption) {
	acknowledgeOrder(button)
	if button.Type == def.CAB {
		light := def.ButtonLight{def.CAB, button.Floor, true}
		buttonLightCh <- light
		if localState.Behaviour == def.IDLE {
			executeOrder(motorDirCh, stuckTimerCh)
		}
	}
}

func executeOrder(motorDirCh chan<- def.MotorDirection, stuckTimerCh chan<- timerOption) {
	dir := elevatorLogic.ChooseDirection(localState)
	switch dir {
	case def.STOP:
		// No order; keep chilling
		localState.Behaviour = def.IDLE

	default:
		localState.Behaviour = def.MOVING
		localState.Direction = dir
		motorDirCh <- dir

		// Elevator about to move, so start the stuck-timer
		stuckTimerCh <- RESET
	}

}

// Clears all hallrequests so they can be reassigned, and sets status to stuck
func elevatorStuck() {
	localState.Stuck = true
	for floor := 0; floor < def.NUMB_FLOORS; floor++ {
		for button := def.HALL_UP; button <= def.HALL_DOWN; button++ {
			localState.AcceptedOrders[floor][button] = false
		}
	}
}

// Close door and turn off light, execute next order
func doorTimerTimedOut(motorDirCh chan<- def.MotorDirection, doorLightCh chan<- bool, stuckTimerCh chan<- timerOption) {
	if localState.Behaviour != def.DOOR_OPEN {
		// In case of errors
		return
	}
	doorLightCh <- false
	executeOrder(motorDirCh, stuckTimerCh)
}
