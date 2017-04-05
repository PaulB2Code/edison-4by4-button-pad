package KeyPadEdison

import (
	"log"
	"os"
	"os/signal"
	"testing"
	"time"
)

var row_out = [4]int{43, 41, 40, 42}   // Give GP Value
var column_in = [4]int{14, 15, 49, 48} // Give GP Value

var keys = [4][4]string{
	{"D", "C", "B", "A"},
	{"#", "9", "6", "3"},
	{"0", "8", "5", "2"},
	{"*", "7", "4", "1"}}

func TestClikedLetter(t *testing.T) {
	log.Println("Start Read Key Value Programm")

	kp, err := New(row_out, column_in, keys)
	if err != nil {
		log.Panic("[ERROR - FATAL], Impossible de démarrer, ", err)
	}

	defer kp.Close()

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt)
	go func() {
		for _ = range closeSignal {
			log.Println("[INFO], Close Systme")
			kp.Close()
			os.Exit(0)
		}
	}()

	c := make(chan int)
	kp.TrackClicked(100*time.Millisecond, c)
	defer kp.StopTracking()

	for {
		kp.Reading = false
		val := <-c

		letter, err := kp.GetValueWithColumn(val)
		if err != nil {
			log.Println("[ERROR], ", err)
		}
		log.Println("Letter clicked", letter)
		//time.Sleep(300 * time.Millisecond)
	}
}
