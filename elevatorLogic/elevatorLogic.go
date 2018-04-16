package elevatorLogic

import (
	"../def"
)

func requestAbove(state def.States) bool {
	for floor := state.Floor + 1; floor < def.NUMB_FLOORS; floor++ {
		for buttonType := 0; buttonType < def.NUMB_BUTTONS; buttonType++ {
			if state.AcceptedOrders[floor][buttonType] == true {
				return true
			}
		}
	}
	return false
}

func requestBelow(state def.States) bool {
	for floor := state.Floor - 1; floor >= 0; floor-- {
		for buttonType := 0; buttonType < def.NUMB_BUTTONS; buttonType++ {
			if state.AcceptedOrders[floor][buttonType] == true {
				return true
			}
		}
	}
	return false
}

//Stop if there is an order in the current direction, cab order or if there are no more orders in the same direction
func ShouldStop(state def.States) bool {
	switch state.Direction {
	case def.DOWN:
		return (state.AcceptedOrders[state.Floor][def.HALL_DOWN] || state.AcceptedOrders[state.Floor][def.CAB] || !requestBelow(state))
	case def.UP:
		return (state.AcceptedOrders[state.Floor][def.HALL_UP] || state.AcceptedOrders[state.Floor][def.CAB] || !requestAbove(state))
	default:
		return true
	}
}

func ChooseDirection(state def.States) def.MotorDirection {
	switch state.Direction {
	case def.UP:
		//If order(s) above, continue upwards. If not, go down or stop.
		if requestAbove(state) {
			return def.UP
		} else if requestBelow(state) {
			return def.DOWN
		} else {
			return def.STOP
		}
	default:
		//Default direction is down, if idle, go down first.
		if requestBelow(state) {
			return def.DOWN
		} else if requestAbove(state) {
			return def.UP
		} else {
			return def.STOP
		}

	}
}

func ClearOrderAtCurrentFloor(state def.States) def.States {
	//They have arrived at their floor and are leaving the elevator
	state.AcceptedOrders[state.Floor][def.CAB] = false

	switch state.Direction {
	case def.UP:
		state = clearOrder(state, def.HALL_UP)
		//if no orders above, the elevator is ready to go down
		if !requestAbove(state) {
			state = clearOrder(state, def.HALL_DOWN)
		}
	case def.DOWN:
		state = clearOrder(state, def.HALL_DOWN)

		//if no orders below, the elevator is ready to go up
		if !requestBelow(state) {
			state = clearOrder(state, def.HALL_UP)

		}
	case def.STOP:

		state = clearOrder(state, def.HALL_UP)
		state = clearOrder(state, def.HALL_DOWN)

	}
	return state
}

// If it has accepted the order, set order to finished and clear from acceptedOrders
func clearOrder(state def.States, buttonType def.ButtonType) def.States {
	if state.AcceptedOrders[state.Floor][buttonType] && buttonType != def.CAB {
		state.HallRequests[state.Floor][buttonType] = 2
		state.AcceptedOrders[state.Floor][buttonType] = false
	} else {
		state.AcceptedOrders[state.Floor][buttonType] = false
	}

	return state
}

// Checks if the elevator should open the door at the current floor (for one particular order)
func ShouldOpenDoor(state def.States, order def.Button) bool {
	return state.Behaviour != def.MOVING && state.Floor == order.Floor &&
		(order.Type != def.HALL_DOWN && requestAbove(state) || order.Type != def.HALL_UP && requestBelow(state) || !requestAbove(state) && !requestBelow(state))
}
