package KeyPadEdison

import (
	"log"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestClikedLetter(t *testing.T) {
	log.Println("Start Read Key Value Programm")

	kp, err := New()
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
