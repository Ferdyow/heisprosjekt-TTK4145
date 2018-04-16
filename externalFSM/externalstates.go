package externalFSM

import (
	"time"

	"../def"
	"../network/peers"
)

//Sorted by ID, contains external states
var externalStateMap = make(map[string]def.States)

func ExternalStateManager(statesToOrderAssignerCh chan<- map[string]def.States,
	externalUpdateCh <-chan def.States, peerUpdateCh <-chan peers.PeerUpdate, localID string) {

	ticker := time.NewTicker(def.STATE_TRANSMIT_INTERVAL)
	for {
		select {
		case <-ticker.C:
			// Create a map to avoid read/write collisions since maps are sent as pointers
			var tempMap = make(map[string]def.States)
			for k, v := range externalStateMap {
				tempMap[k] = v
			}
			statesToOrderAssignerCh <- tempMap
		case peerUpdate := <-peerUpdateCh:
			managePeers(peerUpdate, localID)
		case stateUpdate := <-externalUpdateCh:
			updateExternalState(stateUpdate)
		}
	}
}

// Adds peers that connect to the network and removes peers that disconnect to/from the statemap
func managePeers(peerUpdate peers.PeerUpdate, localID string) {
	if len(peerUpdate.New) != 0 && peerUpdate.New != localID {
		var newState def.States
		newState.Id = peerUpdate.New
		externalStateMap[peerUpdate.New] = newState
	}
	if len(peerUpdate.Lost) != 0 {
		for _, id := range peerUpdate.Lost {
			delete(externalStateMap, id)
		}
	}
}

// Update the states when an update is received
func updateExternalState(update def.States) {
	id := update.Id

	//If the ID is in the stateMap, update it
	if _, ok := externalStateMap[id]; ok {
		externalStateMap[id] = update
	}

}
