package orderAssigner

import (
	"../def"
)

// Return true if an active elevator has accepted the hall request
func orderAccepted(stateMap map[string]def.States, floor int, buttonType def.ButtonType) bool {
	for _, state := range stateMap {
		if state.AcceptedOrders[floor][buttonType] {
			return true
		}
	}
	return false
}

// Returns true if all active elevators have acknowledged the hall request
func orderSynchronized(stateMap map[string]def.States, floor int, buttonType def.ButtonType) bool {
	for _, state := range stateMap {
		if state.HallRequests[floor][buttonType] != 1 {
			return false
		}
	}
	return true
}

// Returns true if the local elevator should accept the hall request
func shouldAcceptOrder(stateMap map[string]def.States, localId string) bool {
	if stateMap[localId].Stuck {
		// Deny orders since the elevator is unable to execute them
		return false
	} else if len(stateMap) == 1 {
		// If alone, take the order
		return true
	}
	// Uses costFunction to find least expensive elevator with respect to time
	var minCost int = 1000000000000000
	var minId string
	for id, state := range stateMap {
		cost := timeToIdle(state)
		if cost < minCost && !state.Stuck {
			minCost = cost
			minId = id

		}
	}
	return minId == localId
}

/************************************
 *	 HALL REQUEST STATUS    		*
 * 0	 -	No order				*
 * 1	 -	order acknowledged		*
 * all 1 -	should accept	order	*
 * 2	 -	order finished			*
 ************************************/
// Runs an algorithm that synchronizes the hall request using the above statuses
func updateLocalOrder(stateMap map[string]def.States, localID string, floor int, buttonType def.ButtonType, buttonLightChan chan<- def.ButtonLight) def.OrderStatus {
	order := def.OrderStatus{buttonType, floor, 0}
	localState := stateMap[localID]

	if hasOrderStatus(stateMap, floor, buttonType, 0) && hasOrderStatus(stateMap, floor, buttonType, 1) && hasOrderStatus(stateMap, floor, buttonType, 2) {
		// Invalid state, retake the order just in case
		if localState.HallRequests[floor][buttonType] != 0 {
			order.Status = 1
		}

	} else if hasOrderStatus(stateMap, floor, buttonType, 2) && hasOrderStatus(stateMap, floor, buttonType, 1) {
		// If at least one has status order finished and one has status acknowledged, update local elevator to order finished
		order.Status = 2
	} else if hasOrderStatus(stateMap, floor, buttonType, 2) {
		// If all have status order finished or no order, clear the order
		if localState.HallRequests[floor][buttonType] != 0 {
			order.Status = 0
			buttonLightChan <- def.ButtonLight{buttonType, floor, false}
		}

	} else if hasOrderStatus(stateMap, floor, buttonType, 1) {
		// If one elevator has acknowledged the order and others have not, set local elevator to acknowledged
		order.Status = 1
	}

	return order
}

// Returns true if one elevator has the given status
func hasOrderStatus(stateMap map[string]def.States, floor int, buttonType def.ButtonType, status int) bool {
	for _, state := range stateMap {
		if state.HallRequests[floor][buttonType] == status {
			return true
		}
	}
	return false
}
