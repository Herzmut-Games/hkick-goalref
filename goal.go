package main

import (
	"fmt"
	"log"
	"net/http"
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
		fmt.Println("White scored")
		client.Publish(goalTopic, 0, false, "white")

		_, err := http.PostForm("http://localhost:3000/goals", url.Values{"team": {"white"}})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
	}
	redGoalLastState = currentState
}

func watchWhiteGoal(pin *gpio.Pin) {
	mutex.Lock()
	defer mutex.Unlock()

	currentState := pin.Read()
	if whiteGoalLastState == true && currentState == false {
		fmt.Println("Red scored")
		client.Publish(goalTopic, 0, false, "red")

		_, err := http.PostForm("http://localhost:3000/goals", url.Values{"team": {"red"}})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)
	}
	whiteGoalLastState = currentState
}

func main() {
	uri, err := url.Parse("mqtt://172.30.1.32:1883")
	if err != nil {
		log.Fatal(err)
	}
	client = connect("pub", uri)
	mutex = &sync.Mutex{}

	err = gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()
	redGoal := gpio.NewPin(17)
	redGoal.Input()
	redGoal.PullUp()
	whiteGoal := gpio.NewPin(27)
	whiteGoal.Input()
	whiteGoal.PullUp()

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	err = redGoal.Watch(gpio.EdgeRising, watchRedGoal)
	if err != nil {
		panic(err)
	}
	defer redGoal.Unwatch()

	err = whiteGoal.Watch(gpio.EdgeRising, watchWhiteGoal)
	if err != nil {
		panic(err)
	}
	defer whiteGoal.Unwatch()

	fmt.Println("Watching Pin goals...")
	select {
	case <-quit:
	}

}
