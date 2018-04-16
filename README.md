# Elevator Project
# Saving of states
Every elevator has a state struct for itself, as well as an externalStateMap that contain the complete states of all participating elevators. The local struct together with the external state map are used together to distribute hall requests. Each elevator can see if another elevator has accepted an order, it is therefore simple to check if an order needs to be reassigned in case elevators lose network or get stuck. 
Hall requests are synchronized using the following table: 
- 0	    -	No order		
- 1	    -	Order acknowledged	
- all 1 -	An elevator should accept order		
- 2	    -	Order finished	
The algorithm increments all of the elevators hall requests to 1 if one elevator has acknowledged an order. When all have acknowledged, one elevator accepts it. When an order has been finished by an elevator, that elevator sets its orderstatus to 2, which the others detect and will update to. When all have status 2, the order is cleared.

# Network interface
The local state is broadcasted with UDP every 50 milliseconds, and it's ID is broadcasted every 15 milliseconds. If an external ID has not been received for 50 milliseconds, this peer is marked as lost. 

# Redundancy
Redundancy is assured by broadcasting the complete states to the network and in addition, cab requests are backed up to a file in case of unexpected crashes. This needs to be done since peers that go offline are deleted from other elevators. 

## External Code
External code used for this project include the network module, except network.go itself and the hardware module which is slightly modified from the version found on the class repository. The cost function has been translated from [1.1 in project resources](https://github.com/TTK4145/Project-resources/tree/master/cost_fns), this uses functions inspired by the [elevator algorithm in D](https://github.com/TTK4145/Project-resources/blob/master/cost_fns/hall_request_assigner/elevator_algorithm.d) which are modified to fit our program. 
