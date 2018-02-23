package localstate

import (
	"../elevio"
	"fmt"
)

const (
	NUMB_FLOOR        = 4
	NUMB_HALL_BUTTONS = 2

	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2

	DOWN = -1
	STOP = 0
	UP   = 1
)

// Define our states

type States struct {
	id           string
	behaviour    int
	floor        int
	direction    int
	HallRequests [NUMB_FLOOR][NUMB_HALL_BUTTONS]int
	CabRequests  [NUMB_FLOOR]int
	isAlive      bool
}

var CurrState States

func UpdateButtonState(button elevio.ButtonEvent) {
	if button.Button == elevio.ButtonType(2) {
		CurrState.CabRequests[button.Floor] = 1
	} else {
		CurrState.HallRequests[button.Floor][button.Button] = 1
	}
	fmt.Println("hallRequessts: ", CurrState.HallRequests)
	fmt.Println("CabRequests: ", CurrState.CabRequests)
}
