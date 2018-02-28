package main

import (
	"./elevio"
	"./states"
	"./network"
	"./network/bcast"
	//"./network/localip"
	"./network/peers"
	"fmt"
)

func main() {
	fmt.Println("Hello")

	//Initialize the hardware and ID
	id := network.Init()
	states.Init(id)
	elevio.Init("localhost:15657", states.NUMB_FLOOR)

	var dir elevio.MotorDirection = elevio.MD_Up
	fmt.Println("direction: ", dir)

	//Make a channel to send states to the network
	statesToNetwork := make(chan states.States)
	
	go states.SendStatesOnInterval(statesToNetwork) // Puts states on the channel 
	go bcast.Transmitter(20014, statesToNetwork) //send states when they are received on channel



	ButtonPressedCh := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(ButtonPressedCh)


	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(30014, peerUpdateCh)
	go states.ManagePeers(peerUpdateCh) //Update peers when an update is received

	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTransmitEnable := make(chan bool)
	go peers.Transmitter(30014, id, peerTransmitEnable)
	
	

	

	// We make channels for sending and receiving our custom data types
	stateRecCh := make(chan states.States)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(20014, statesToNetwork)
	go bcast.Receiver(20014, stateRecCh)

	go states.UpdateExternalState(stateRecCh)

	for {
		select {
		case button := <-ButtonPressedCh:
			states.UpdateButtonState(button)
			//stateTransCh <- states.CurrState
			elevio.SetButtonLamp(button.Button, button.Floor, true) //don't do this here later

		//case externalState := <-stateRecCh:
			//fmt.Println("RECEIVED UPDATE: ", externalState)
		}
	}
}
