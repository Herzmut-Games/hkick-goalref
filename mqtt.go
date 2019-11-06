package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	client    mqtt.Client
	goalTopic = "goals"
)

func connect(clientID string, uri *url.URL) mqtt.Client {
	uri, err := url.Parse("mqtt://172.30.1.32:1883")
	if err != nil {
		log.Fatal(err)
	}
	
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func goalCount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	scoredTeam := r.FormValue("team")

	fmt.Printf("Team %s scored", scoredTeam)
	client.Publish(goalTopic, 0, false, scoredTeam)
}
