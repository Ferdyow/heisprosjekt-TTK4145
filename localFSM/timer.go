package localFSM

import (
	"time"

	"../def"
)

type timerOption int

const (
	RESET timerOption = 0
	STOP              = 1
)

func doorTimer(doorTimeOutCh chan<- bool, doorResetCh <-chan bool) {
	// Initialize the timer
	doorTimer := time.NewTimer(def.DOOR_OPEN_DURATION)
	doorTimer.Stop()
	for {
		select {
		case <-doorTimer.C:
			doorTimeOutCh <- true

		case <-doorResetCh:
			doorTimer.Reset(def.DOOR_OPEN_DURATION)
		}
	}
}

func stuckTimer(isStuckCh chan<- bool, stuckTimerCh <-chan timerOption) {
	// Initialize the timer
	stuckTimer := time.NewTimer(def.ELEVATOR_STUCK_DURATION)
	stuckTimer.Stop()
	for {
		select {
		case <-stuckTimer.C:
			isStuckCh <- true

		case state := <-stuckTimerCh:
			switch state {
			case RESET:
				stuckTimer.Reset(def.ELEVATOR_STUCK_DURATION)
			case STOP:
				stuckTimer.Stop()
			}
		}
	}
}
