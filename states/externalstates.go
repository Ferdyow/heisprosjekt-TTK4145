package states

import (
	//"../elevio"
	"fmt"
	"time"
	"../network/peers"
)


//Sorted by ID, contains external states
var externalStateMap = make(map[string]States)




//Function that adds and initializes a peer
func ManagePeers(PeerUpdateCh <-chan peers.PeerUpdate){
	PeerUpdate := <- PeerUpdateCh
	if  len(PeerUpdate.New) != 0 && PeerUpdate.New != localState.id{
		var newState States
		externalStateMap[PeerUpdate.New] = newState;
	}
	if len(PeerUpdate.Lost) != 0 {
		for _,id := range PeerUpdate.Lost{
			delete(externalStateMap, id)
		}
	}
	fmt.Println("HEEEEEEEEEEEEEEEEEEEELLLLLLLLLLLLLLLLLLLLLOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
	fmt.Println(externalStateMap)
	time.Sleep(1*time.Second)

}