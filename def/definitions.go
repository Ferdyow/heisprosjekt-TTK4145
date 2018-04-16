package def

import "time"

const (
	NUMB_FLOORS       = 4
	NUMB_HALL_BUTTONS = 2
	NUMB_BUTTONS      = 3

	//Ports used for UDP
	BroadcastPort = 20014
	PeersPort     = 30014

	//Ports used for TCP
	ServerPort = "15657"
	SimPort1   = "15555"
	SimPort2   = "15556"
	SimPort3   = "15557"
)

const ELEVATOR_STUCK_DURATION = time.Second * 6
const DOOR_OPEN_DURATION = time.Second * 3
const STATE_TRANSMIT_INTERVAL = time.Millisecond * 50

/************************************
 *	 HALL REQUEST STATUS    		*
 * 0	 -	No order				*
 * 1	 -	order acknowledged		*
 * all 1 -	should accept order		*
 * 2	 -	order finished			*
 ************************************/

// Complete state of one elevator
type States struct {
	Id             string
	Behaviour      Behaviour
	Floor          int
	Direction      MotorDirection
	HallRequests   [NUMB_FLOORS][NUMB_HALL_BUTTONS]int
	AcceptedOrders [NUMB_FLOORS][NUMB_BUTTONS]bool
	Stuck          bool
}

type Behaviour int

const (
	IDLE      Behaviour = 0
	MOVING              = 1
	DOOR_OPEN           = 2
)

type MotorDirection int

const (
	UP   MotorDirection = 1
	DOWN                = -1
	STOP                = 0
)

type ButtonType int

const (
	HALL_UP   ButtonType = 0
	HALL_DOWN            = 1
	CAB                  = 2
)

type Button struct {
	Floor int
	Type  ButtonType
}

type OrderStatus struct {
	Type   ButtonType
	Floor  int
	Status int
}

type ButtonLight struct {
	Type   ButtonType
	Floor  int
	Status bool
}
