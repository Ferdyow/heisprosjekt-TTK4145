package network

import (
	"flag"
	"fmt"
	"os"
	//"time"

	//"./network/bcast"
	"./localip"
	//"./peers"
)

var id string

func Init() string {
	SetId()
	fmt.Println("Local id: ", id)
	return id
}

func SetId() {
	//Set the elevator ID in terminal
	//  `go run main.go -id "our_id"`

	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
}


