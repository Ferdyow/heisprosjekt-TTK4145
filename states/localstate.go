package states

import (
	"../elevio"
	"fmt"
	"time"
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
// should probably move this thing
type States struct {
	Id           string 
	Behaviour    int
	Floor        int
	Direction    int
	HallRequests [NUMB_FLOOR][NUMB_HALL_BUTTONS]int
	CabRequests  [NUMB_FLOOR]int
	IsAlive      bool
}

var LocalState States

func Init(id string){
	LocalState.Id = id;
	LocalState.Behaviour = MOVING;
	LocalState.Direction = UP;
	//fmt.Println("LocalState: ", LocalState)
}





func UpdateButtonState(button elevio.ButtonEvent) {
	if button.Button == elevio.ButtonType(2) {
		LocalState.CabRequests[button.Floor] = 1
	} else {
		LocalState.HallRequests[button.Floor][button.Button] = 1
	}
	//fmt.Println("hallRequessts: ", LocalState.HallRequests)
	//fmt.Println("CabRequests: ", LocalState.CabRequests)
}

func SendStatesOnInterval(statesToNetworkChan chan<- States){
	tick := time.NewTicker(200*time.Millisecond)
	for{
		select{
		case <- tick.C:
			fmt.Println("Localstate:", LocalState)
			statesToNetworkChan <- LocalState
		}
	}	
}