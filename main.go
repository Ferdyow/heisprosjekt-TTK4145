package main

import (
	"time"

	"./backup"
	"./def"
	"./externalFSM"
	"./hardware"
	"./localFSM"
	"./network"
	"./network/peers"
	"./orderAssigner"
)

// Sets the hardware to a valid position, gives the elevator an ID and sends the position to localFSM
func initialize(floorSensorCh chan<- int) (string, int) {
	// initialize local elevator
	//Initialize the hardware and ID
	id := network.Init()

	// Switch between:
	// def.ServerPort
	// def.SimPort1,2,3
	initialFloor := hardware.Init("localhost:"+string(def.ServerPort), def.NUMB_FLOORS)
	return id, initialFloor
}

func main() {
	statesToBackupCh := make(chan def.States)

	/********************************NETWORK CHANNELS********************************/
	// Channels for sending and receiving elevator states
	statesToNetworkCh := make(chan def.States)
	stateRecCh := make(chan def.States)

	// Channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)

	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable"
	peerTransmitEnable := make(chan bool)
	/*******************************************************************************/

	/*******************************HARDWARE CHANNELS*******************************/
	buttonPressedCh := make(chan def.Button)
	floorSensorCh := make(chan int)
	buttonLightCh := make(chan def.ButtonLight)
	motorDirectionCh := make(chan def.MotorDirection)
	doorLightCh := make(chan bool)
	/*******************************************************************************/

	// Channel that sends hall requests to be accepted by the local elevator from the order assigner to localFSM
	acceptOrderCh := make(chan def.Button)

	// Channel that sends the status of hall requests to localFSM from the order assigner
	hallReqCh := make(chan def.OrderStatus)

	// Channels that send the local and external states to the order assigner
	statesToOrderAssignerCh := make(chan def.States)
	externalStatesCh := make(chan map[string]def.States)

	id, initialFloor := initialize(floorSensorCh)

	// Each module runs its own eventmanager, the only thing done in main is creating channels
	// that allow communication between different modules
	go backup.BackupManager(statesToBackupCh, buttonPressedCh)
	go externalFSM.ExternalStateManager(externalStatesCh, stateRecCh, peerUpdateCh, id)
	go network.NetworkManager(statesToNetworkCh, stateRecCh, peerUpdateCh, peerTransmitEnable)
	go hardware.HardwareManager(buttonPressedCh, floorSensorCh, buttonLightCh, motorDirectionCh, doorLightCh)
	go orderAssigner.OrderManager(statesToOrderAssignerCh, externalStatesCh, hallReqCh, acceptOrderCh, buttonLightCh)

	go localFSM.EventManager(hallReqCh, acceptOrderCh,
		statesToNetworkCh, statesToBackupCh, statesToOrderAssignerCh,
		floorSensorCh, buttonPressedCh, buttonLightCh,
		motorDirectionCh, doorLightCh, id, initialFloor)

	for {
		time.Sleep(1 * time.Second)
	}

}
