package localFSM

import (
	//"../elevio"
	//"fmt"
	//"time"
	//"../network/peers"
	"../states"
	"../def"
)


//helper function
func RequestAbove(state states.States) bool {
	for floor := state.Floor; floor < def.NUMB_FLOOR; floor++ {
		for button := 0; button < def.NUMB_BUTTONS; button++ {
			if state.AcceptedOrders[floor][button] == true {
				return true
			}
		} 
	} 
	return false
}

//helper function
func RequestBelow(state states.States) bool {
	for floor := 0; floor < state.Floor; floor++ {
		for button := 0; button < def.NUMB_BUTTONS; button++ {
			if state.AcceptedOrders[floor][button] == true {
				return true
			}
		} 
	} 
	return false
}


func ShouldStop(state states.States) bool{
	switch state.Direction {
	case def.DOWN:
		return (state.AcceptedOrders[state.Floor][def.BUTTON_DOWN] || state.AcceptedOrders[state.Floor][def.BUTTON_CAB]  || !RequestBelow(state));
	case def.UP:
		return (state.AcceptedOrders[state.Floor][def.BUTTON_UP] || state.AcceptedOrders[state.Floor][def.BUTTON_CAB]  || !RequestAbove(state));
	default:
		return true;
	}
}


func chooseDir(state states.States) int{
	switch state.Direction{
	case def.UP:
		//if order above, continue up. if not go down or stop
		if RequestAbove(state){
			return def.UP;
		} else if RequestBelow(state){
			return def.DOWN;
		} else{
			return def.STOP;
		}
	default:
	//Default direction is down, if standing still, go down first
		if RequestBelow(state){
			return def.DOWN;
		} else if RequestAbove(state){
			return def.UP;
		} else{
			return def.STOP;
		}
		
	}
}


func clearOrderAtCurrentFloor(state states.States) states.States {
	//They have arrived at their floor and are leaving the elevator
	state.AcceptedOrders[state.Floor][def.BUTTON_CAB] = false;
	
	switch state.Direction{
	case def.UP:
		state.AcceptedOrders[state.Floor][def.BUTTON_UP] = false;
		//if no orders above, the elevator is ready to go down
		if !RequestAbove(state){
			state.AcceptedOrders[state.Floor][def.BUTTON_DOWN] = false;
		}
	case def.DOWN:
		state.AcceptedOrders[state.Floor][def.BUTTON_DOWN] = false;
		//if no orders below, the elevator is ready to go up
		if !RequestBelow(state){
			state.AcceptedOrders[state.Floor][def.BUTTON_UP] = false;
		}
	case def.STOP:
		state.AcceptedOrders[state.Floor][def.BUTTON_UP] = false;
		state.AcceptedOrders[state.Floor][def.BUTTON_DOWN] = false;

	}
	return state;
}