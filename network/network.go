package network

import (
	"flag"
	"fmt"
	"os"

	"../def"
	"./bcast"
	"./localip"
	"./peers"
)

var id string

func Init() string {
	setId()
	fmt.Println("Local id: ", id)
	return id
}

func setId() {
	//Set the elevator ID in terminal
	//  `go run main.go -id "our_id"`
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address and process ID
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
}

// Manage receiving and broadcasting of both elevator states and peer updates
func NetworkManager(statesToNetworkCh <-chan def.States, stateRecCh chan<- def.States, peerUpdateCh chan<- peers.PeerUpdate, peerTransmitEnable <-chan bool) {
	go bcast.Transmitter(def.BroadcastPort, statesToNetworkCh)
	go bcast.Receiver(def.BroadcastPort, stateRecCh)
	go peers.Receiver(def.PeersPort, peerUpdateCh)
	go peers.Transmitter(def.PeersPort, id, peerTransmitEnable)
}
