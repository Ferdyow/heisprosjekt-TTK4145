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

//set the states
func UpdateExternalState(externalUpdateCh <-chan States){
	var id string
	var update States
	for{
		update = <- externalUpdateCh
		id = update.Id
		fmt.Println(update)
		fmt.Println("ID: ", id)
		if _, ok := externalStateMap[id]; ok {
			externalStateMap[id] = update
			fmt.Println("external states updated: ", externalStateMap)
		}
	}
}