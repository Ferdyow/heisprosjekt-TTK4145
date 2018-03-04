package localFSM

import (
	//"../elevio"
	//"fmt"
	//"time"
	//"../network/peers"
	"../states"
	"../def"
)



func handleOrderAssigned(buttonCh <-chan def.Button, state states.States){
	for{ //FOREVER
		button := <-buttonCh;
		switch state.Behaviour{
		case def.DOOR_OPEN:
			if state.Floor == button.Floor{
				//Restart the timer
			} else {
				//set the fucking order
			}
			case 

		}
	}
}