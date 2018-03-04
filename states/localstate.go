package states

import (
	"../elevio"
	"fmt"
	"time"
	"../def"
)



// Define our states
// should probably move this thing
type States struct {
	Id           string 
	Behaviour    int
	Floor        int
	Direction    int
	HallRequests [def.NUMB_FLOOR][def.NUMB_HALL_BUTTONS]int
	CabRequests  [def.NUMB_FLOOR]int
	AcceptedOrders [def.NUMB_FLOOR][def.NUMB_BUTTONS]bool
	IsAlive      bool
}

var LocalState States

func Init(id string){
	LocalState.Id = id;
	LocalState.Behaviour = def.MOVING;
	LocalState.Direction = def.UP;
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
			fmt.Println("Localstate:\t\t", LocalState)
			statesToNetworkChan <- LocalState
		}
	}	
}