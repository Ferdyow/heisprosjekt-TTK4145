package backup

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"../def"
)

// Puts recovered orders on buttonPressedCh and backs up states when they are put on the channel
func BackupManager(AllStatesCh <-chan def.States, buttonPressedCh chan<- def.Button) {
	Orders := recoverCabOrders()
	for floor := 0; floor < def.NUMB_FLOORS; floor++ {
		if Orders[floor] == 1 {
			buttonPressedCh <- def.Button{floor, def.CAB}
		}
	}
	backupCabOrders(AllStatesCh)
}

// Saves the cabRequests to the file backup.txt in the same directory as main.go
func backupCabOrders(AllStatesCh <-chan def.States) {
	var state def.States
	var cabRequests [def.NUMB_FLOORS]int

	for {
		// Wait for state to be received
		state = <-AllStatesCh

		// Create backup file in same directory
		backupFile, err := os.Create("backup.txt")
		if err != nil {
			log.Fatal(err)
		}

		//Make an array of integers to store
		for floor := 0; floor < def.NUMB_FLOORS; floor++ {
			if state.AcceptedOrders[floor][def.CAB] {
				cabRequests[floor] = 1
			} else {
				cabRequests[floor] = 0
			}
		}

		//Convert int array to string array and write it to buffer
		ToString := strings.Fields(strings.Trim(fmt.Sprint(cabRequests), "[]"))
		WriteBackup := csv.NewWriter(backupFile)
		err = WriteBackup.Write(ToString)

		if err != nil {
			log.Fatal(err)
		}
		WriteBackup.Flush()
		backupFile.Close()
	}
}

// Recovers and returns the cabRequests saved in backup.txt in the same directory as main.go
//as an array of integers (1 = true, 0 = false)
func recoverCabOrders() [def.NUMB_FLOORS]int {
	var cabRequests [def.NUMB_FLOORS]int
	backupFile, err := ioutil.ReadFile("backup.txt")

	if err != nil {
		return cabRequests
	}
	ReadBackup := csv.NewReader(strings.NewReader(string(backupFile)))
	stringArray := []string{}
	stringArray, err = ReadBackup.Read()

	if err == io.EOF {
		return cabRequests
	}
	for i := 0; i < def.NUMB_FLOORS; i++ {
		cabRequests[i], _ = strconv.Atoi(stringArray[i])
	}
	return cabRequests
}
