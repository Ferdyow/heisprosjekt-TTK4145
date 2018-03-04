package def

const (
	NUMB_FLOOR        = 4
	NUMB_HALL_BUTTONS = 2
	NUMB_BUTTONS = 3

	//Behaviour
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2

	//Direction
	DOWN = -1
	STOP = 0
	UP   = 1

	//ButtonType
	BUTTON_DOWN = 0
	BUTTON_UP = 1
	BUTTON_CAB = 2
)

type  Button struct {
	Floor int;
	Buttontype int;
}