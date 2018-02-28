package states

import (
	//"../elevio"
	"fmt"
	//"time"
	"../network/peers"
)


//Sorted by ID, contains external states
var externalStateMap map[string]States
//externalStateMap := make(map[string]states)

//Function that adds and initializes a peer
func ManagePeers(PeerUpdateCh <-chan peers.PeerUpdate){
	PeerUpdate := <- PeerUpdateCh
	if  PeerUpdate.New != "" && PeerUpdate.New != localState.id{
		var newState States
		externalStateMap[PeerUpdate.New] = newState;
	}
	if len(PeerUpdate.Lost) != 0 {
		for _,id := range PeerUpdate.Lost{
			delete(externalStateMap, id)
		}
	}
	fmt.Println(externalStateMap)

}