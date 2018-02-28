package states

import (
	//"../elevio"
	"fmt"
	//"time"
	"../network/peers"
)


//Sorted by ID, contains external states
var externalStateMap = make(map[string]States)




//Function that adds and initializes a peer
func ManagePeers(PeerUpdateCh <-chan peers.PeerUpdate){
	for{
		PeerUpdate := <- PeerUpdateCh
		if  len(PeerUpdate.New) != 0 && PeerUpdate.New != LocalState.Id{
			var newState States
			externalStateMap[PeerUpdate.New] = newState;
		}
		if len(PeerUpdate.Lost) != 0 {
			for _,id := range PeerUpdate.Lost{
				delete(externalStateMap, id)
			}
		}
		fmt.Println("Peers updated: ", externalStateMap)
	}
}

//set the states | there are some timing issues here |
func UpdateExternalState(externalUpdateCh <-chan States){
	var id string
	var update States
	for{
		update = <- externalUpdateCh
		id = update.Id
		if _, ok := externalStateMap[id]; ok {
			externalStateMap[id] = update
			for _, value := range externalStateMap{
				fmt.Println("external states:\t", value, "\n\n")
		     }
		}
	}
}