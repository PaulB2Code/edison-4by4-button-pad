package KeyPadEdison

import (
	"errors"
	"fmt"
	"sync"
	"time"

	edio "github.com/PaulB2Code/edison-gpio"
)

var row_out_default = [4]int{43, 41, 129, 128} // Give GP Value
var column_in_default = [4]int{40, 42, 12, 13} // Give GP Value

//Give the Key of the Key Pad
var keys_default = [4][4]string{
	{"1", "2", "3", "A"},
	{"4", "5", "6", "B"},
	{"7", "8", "9", "C"},
	{"*", "0", "#", "D"}}

var colPin = make([]int, 0, 0)
var rowPin = make([]int, 0, 0)

func (kp *Keypad) Close() {
	// Unmap gpio memory when done
}

type Keypad struct {
	Reading   bool
	Lock      sync.RWMutex
	row_out   [4]int
	column_in [4]int
	keys      [4][4]string
}

func New(row_out [4]int, column_in [4]int, keys [4][4]string) (Keypad, error) {
	//Index Ouput to High Value
	for _, val := range row_out {
		rowPin = append(rowPin, val)
		edio.ExportPin(val)
		edio.ModePin("0", val)
		edio.DirectionPin("out", val)
		edio.ValuePin("1", val)
	}
	for _, val := range column_in {
		colPin = append(colPin, val)
		edio.ExportPin(val)
		edio.ModePin("0", val)
		edio.DirectionPin("in", val)

		//colPin[i].PullDown()
	}
	return Keypad{row_out: row_out, column_in: column_in, keys: keys}, nil
}

func (kp *Keypad) getLetter(row int, col int) (string, error) {
	if len(kp.keys) < row || len(kp.keys) < col {
		return "", errors.New(fmt.Sprintf("No data for row %v and col %v", row, col))
	}
	val := kp.keys[row][col]

	return val, nil
}

var ticker *time.Ticker

func (kp *Keypad) TrackClicked(timeScanning time.Duration, c chan int) {
	ticker = time.NewTicker(timeScanning)
	go func() {
		for _ = range ticker.C {
			if !kp.Reading {
				for i, val := range colPin {
					valPin, _ := edio.ReadPinState(val)
					if valPin == 1 {
						if len(c) == 0 {
							c <- i
						}
						goto end
					}
				}
			}
		end:
		}
	}()

}
func (kp *Keypad) StopTracking() {
	ticker.Stop()
}
func (kp *Keypad) restartTracking(timeScanning time.Duration) {
	ticker = time.NewTicker(timeScanning)
}

func (kp *Keypad) GetValueWithColumn(i int) (string, error) {
	kp.Reading = true
	kp.Lock.RLock()
	defer kp.Lock.RUnlock()
	var rowId int
	//Test to know whos button is pushed
	/*log.Println("[DEBUG], Put to 0 All pin")
	for ii := 1; ii < len(rowPin); ii++ {
		edio.ValuePin("0", rowPin[ii])
	}
	valPin, _ := edio.ReadPinState(colPin[i])
	if valPin == 1 {
		rowId = 0
		goto endblock
	}
	*/
	//log.Println("[DEBUG], Put to 1 Pin by PIn")
	for ii := 0; ii < len(rowPin); ii++ {
		//time.Sleep(10 * time.Millisecond)
		edio.ValuePin("0", rowPin[ii])
		//time.Sleep(5 * time.Millisecond)
		valPin, _ := edio.ReadPinState(colPin[i])
		if valPin == 0 {
			rowId = ii
			goto endblock
		}
	}
	for ii := 0; ii < len(rowPin); ii++ {
		edio.ValuePin("1", rowPin[ii])
	}

	return "", errors.New("Error Reading second number")
endblock:

	letter, err := kp.getLetter(rowId, i)
	if err != nil {
		return "", err
	}
	//Set High All
	//time.Sleep(300 * time.Millisecond)
	for ii := 0; ii < rowId+1; ii++ {
		edio.ValuePin("1", rowPin[ii])
	}
	return letter, nil
}
