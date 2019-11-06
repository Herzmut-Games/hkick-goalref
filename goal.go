package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/warthog618/gpio"
)

var (
	mutex              *sync.Mutex
	redGoalLastState   gpio.Level
	whiteGoalLastState gpio.Level
)

func watchRedGoal(pin *gpio.Pin) {
	mutex.Lock()
	defer mutex.Unlock()

	currentState := pin.Read()
	if redGoalLastState == true && currentState == false {
		// if red goal callback is triggered team white scored
		fmt.Println("White scored")
		client.Publish(goalTopic, 0, false, "white")

		time.Sleep(2 * time.Second)
	}
	redGoalLastState = currentState
}

func watchWhiteGoal(pin *gpio.Pin) {
	mutex.Lock()
	defer mutex.Unlock()

	currentState := pin.Read()
	if whiteGoalLastState == true && currentState == false {
		// if  white goal callback is triggered team red scored
		fmt.Println("Red scored")
		client.Publish(goalTopic, 0, false, "red")
		time.Sleep(2 * time.Second)
	}
	whiteGoalLastState = currentState
}

func main() {
	// connect to mqtt broker as publisher
	uri, err := url.Parse("mqtt://172.30.1.32:1883")
	if err != nil {
		log.Fatal(err)
	}
	client = connect("pub", uri)

	mutex = &sync.Mutex{}

	// try to open gpio pins
	err = gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	// define pin for team red and configure it
	redGoal := gpio.NewPin(17)
	redGoal.Input()
	redGoal.PullUp()

	// define pin for team white and configure it
	whiteGoal := gpio.NewPin(27)
	whiteGoal.Input()
	whiteGoal.PullUp()

	// set callback for red goal
	err = redGoal.Watch(gpio.EdgeRising, watchRedGoal)
	if err != nil {
		panic(err)
	}
	defer redGoal.Unwatch()

	// set callback for white goal
	err = whiteGoal.Watch(gpio.EdgeRising, watchWhiteGoal)
	if err != nil {
		panic(err)
	}
	defer whiteGoal.Unwatch()

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	fmt.Println("Watching Pin goals...")
	select {
	case <-quit:
	}

}
