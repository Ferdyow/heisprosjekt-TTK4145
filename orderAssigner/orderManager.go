package orderAssigner

import (
	"../def"
)

// Uses the states of all active elevators to assign hallrequests
func OrderManager(localStateCh <-chan def.States, externalStatesCh <-chan map[string]def.States, sendOrderStatusCh chan<- def.OrderStatus,
	acceptOrderCh chan<- def.Button, buttonLightCh chan<- def.ButtonLight) {
	for {
		// Wait for the state updates to be received
		stateMap := <-externalStatesCh
		localState := <-localStateCh

		localId := localState.Id
		stateMap[localId] = localState

		for floor := 0; floor < def.NUMB_FLOORS; floor++ {
			for buttonType := def.HALL_UP; buttonType <= def.HALL_DOWN; buttonType++ {
				if orderSynchronized(stateMap, floor, buttonType) {
					// All elevators have acknowledged the order
					buttonLightCh <- def.ButtonLight{buttonType, floor, true}
					if !orderAccepted(stateMap, floor, buttonType) && shouldAcceptOrder(stateMap, localId) {
						// Local elevator accepts the order
						acceptOrderCh <- def.Button{floor, buttonType}
					}
				} else {
					// Run order-synchronization algorithm
					order := updateLocalOrder(stateMap, localId, floor, buttonType, buttonLightCh)
					sendOrderStatusCh <- order
				}

			}
		}
	}
}
