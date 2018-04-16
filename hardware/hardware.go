package hardware

import "time"
import "sync"
import "net"
import "fmt"
import "../def"

const _pollRate = 20 * time.Millisecond

var _initialized bool = false
var _numFloors int = def.NUMB_FLOORS
var _mtx sync.Mutex
var _conn net.Conn

// Initialize to a valid state with all lights off
func Init(addr string, numFloors int) int {
	if _initialized {
		fmt.Println("Driver already initialized!")
		return -1
	}
	_numFloors = numFloors
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	for floor := 0; floor < _numFloors; floor++ {
		for button := def.ButtonType(0); button < def.NUMB_BUTTONS; button++ {
			setButtonLamp(button, floor, false)
		}
	}
	setDoorOpenLamp(false)

	for getFloor() == -1 {
		setMotorDirection(def.DOWN)
	}
	setMotorDirection(def.STOP)

	_initialized = true
	return getFloor()
}

func setMotorDirection(dir def.MotorDirection) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{1, byte(dir), 0, 0})
}

func setButtonLamp(button def.ButtonType, floor int, value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{2, byte(button), byte(floor), toByte(value)})
}

func setFloorIndicator(floor int) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{3, byte(floor), 0, 0})
}

func setDoorOpenLamp(value bool) {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{4, toByte(value), 0, 0})
}

func PollButtons(receiver chan<- def.Button) {
	prev := make([][3]bool, _numFloors)
	for {
		time.Sleep(_pollRate)
		for f := 0; f < _numFloors; f++ {
			for b := def.ButtonType(0); b < 3; b++ {
				v := getButton(b, f)
				if v != prev[f][b] && v != false {
					receiver <- def.Button{f, def.ButtonType(b)}
				}
				prev[f][b] = v
			}
		}
	}
}

func PollFloorSensor(receiver chan<- int, floorIndicatorCh chan<- int) {
	prev := -1
	for {
		time.Sleep(_pollRate)
		v := getFloor()
		if v != prev && v != -1 {
			receiver <- v
			floorIndicatorCh <- v
		}
		prev = v
	}
}

func getButton(buttonType def.ButtonType, floor int) bool {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{6, byte(buttonType), byte(floor), 0})
	var buf [4]byte
	_conn.Read(buf[:])
	return toBool(buf[1])
}

func getFloor() int {
	_mtx.Lock()
	defer _mtx.Unlock()
	_conn.Write([]byte{7, 0, 0, 0})
	var buf [4]byte
	_conn.Read(buf[:])
	if buf[1] != 0 {
		return int(buf[2])
	} else {
		return -1
	}
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 {
		b = true
	}
	return b
}

func HardwareManager(buttonPressedCh chan<- def.Button, floorSensorCh chan<- int, buttonLightCh <-chan def.ButtonLight,
	motorDirectionCh <-chan def.MotorDirection, doorLightCh <-chan bool) {

	floorIndicatorCh := make(chan int)
	go PollFloorSensor(floorSensorCh, floorIndicatorCh)
	go PollButtons(buttonPressedCh)

	for {
		select {
		case buttonLight := <-buttonLightCh:
			setButtonLamp(buttonLight.Type, buttonLight.Floor, buttonLight.Status)
		case direction := <-motorDirectionCh:
			setMotorDirection(direction)
		case doorLight := <-doorLightCh:
			setDoorOpenLamp(doorLight)
		case floorIndicator := <-floorIndicatorCh:
			setFloorIndicator(floorIndicator)
		}
	}
}
