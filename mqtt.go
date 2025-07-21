package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var counter int

var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("[RECV] %s | %s\n", msg.Topic(), string(msg.Payload()))

	response := fmt.Sprintf("Result message #%d", counter)
	counter++

	token := client.Publish("test/result", 0, false, response)
	token.Wait()
	if token.Error() != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Publish failed: %v\n", token.Error())
	} else {
		fmt.Printf("[SEND] test/result => %s\n", response)
	}
}

func main() {
	fmt.Println("[BOOT] Starting MQTT routing server...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	broker := "tcp://broker.hivemq.com:1883"
	clientID := "golang-mqtt-backend"
	topic := "test4472/#"

	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetDefaultPublishHandler(messageHandler)

	opts.OnConnect = func(c MQTT.Client) {
		fmt.Println("[INFO] Connected to broker.")
		if token := c.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Subscribe failed: %v\n", token.Error())
		} else {
			fmt.Printf("[INFO] Subscribed to topic: %s\n", topic)
		}
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Connect failed: %v\n", token.Error())
		os.Exit(1)
	}

	<-sigs
	client.Disconnect(250)
	fmt.Println("[INFO] Disconnected cleanly.")
}

