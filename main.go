package main

import (
	"./elevio"
	"./localstate"
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
	elevio.Init("localhost:15657", localstate.NUMB_FLOOR)

	var dir elevio.MotorDirection = elevio.MD_Up
	fmt.Println(dir)

	//Set up channels
	ButtonPressedCh := make(chan elevio.ButtonEvent)
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)
	go elevio.PollButtons(ButtonPressedCh)

	// We make channels for sending and receiving our custom data types
	stateTransCh := make(chan localstate.States)
	stateRecCh := make(chan localstate.States)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, stateTransCh)
	go bcast.Receiver(16569, stateRecCh)

	for {
		select {
		case button := <-ButtonPressedCh:
			localstate.UpdateButtonState(button)
			stateTransCh <- localstate.CurrState
			elevio.SetButtonLamp(button.Button, button.Floor, true) //don't do this here later

		case externalState := <-stateRecCh:
			fmt.Println(externalState.HallRequests)
		}
	}
}
