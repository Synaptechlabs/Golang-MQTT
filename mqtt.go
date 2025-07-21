package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var counter int

// MQTT message handler
var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("[RECV] %s | %s\n", msg.Topic(), msg.Payload())
	response := fmt.Sprintf("Result message #%d", counter)
	counter++

	if client.IsConnected() {
		pubToken := client.Publish("nn/result", 0, false, response)
		pubToken.Wait()
		if pubToken.Error() != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Publish failed: %v\n", pubToken.Error())
		} else {
			fmt.Printf("[SEND] nn/result => %s\n", response)
		}
	} else {
		fmt.Fprintln(os.Stderr, "[WARN] MQTT client disconnected. Cannot publish.")
	}
}

func main() {
	fmt.Println("[BOOT] Starting MQTT routing server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\n[INFO] Shutdown signal caught. Cleaning up...")
		cancel()
	}()

	broker := "tcp://broker.hivemq.com:1883"
	clientID := "synaptech-wildcard"
	topic := "#"

	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetDefaultPublishHandler(messageHandler)

	opts.OnConnect = func(c MQTT.Client) {
	fmt.Printf("[INFO] Connected to broker: %s\n", broker)

	token := c.Subscribe("#", 0, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[RECV] %s | %s\n", msg.Topic(), msg.Payload())
	})
	token.Wait()
	if token.Error() != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Subscribe failed: %v\n", token.Error())
	} else {
		fmt.Printf("[INFO] Subscribed to topic: %s\n", topic)
	}
}

	client := MQTT.NewClient(opts)
	connToken := client.Connect()
	connToken.Wait()
	if connToken.Error() != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Initial connect failed: %v\n", connToken.Error())
		os.Exit(1)
	}

	<-ctx.Done()
	client.Disconnect(250)
	fmt.Println("[INFO] Disconnected cleanly.")
}
