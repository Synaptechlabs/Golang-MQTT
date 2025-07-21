// Golang MQTT Routing and Forwarding Server
// Author: Scott Douglass (SynaptechLabs.ai)
// Description: Lightweight MQTT listener and re-publisher with reconnect logic,
// auto-recovery, and filtered topic handling to prevent response loops.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var counter int

// MQTT message handler
var messageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())
	fmt.Printf("[RECV] %s | %s\n", topic, payload)

	// Avoid responding to messages on result channels to prevent feedback loops
	if strings.HasSuffix(topic, "/result") {
		fmt.Println("[INFO] Ignoring response message to prevent loop.")
		return
	}

	response := fmt.Sprintf("Result message #%d %s", counter, payload)
	counter++

	// Define output topic
	outputTopic := topic + "/result"

	if client.IsConnected() {
		token := client.Publish(outputTopic, 0, false, response)
		token.Wait()
		if token.Error() != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Publish failed: %v\n", token.Error())
		} else {
			fmt.Printf("[SEND] %s => %s\n", outputTopic, response)
		}
	} else {
		fmt.Fprintln(os.Stderr, "[WARN] MQTT client disconnected. Cannot publish.")
	}
}

func main() {
	fmt.Println("[BOOT] Starting MQTT routing server...")

	// Handle CTRL+C shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\n[INFO] Shutdown signal caught. Cleaning up...")
		cancel()
	}()

	// MQTT connection settings
	broker := "tcp://broker.hivemq.com:1883"
	clientID := "synaptech-wildcard"
	subTopic := "test4472/#" // listen to test4472 and any subtopics

	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetDefaultPublishHandler(messageHandler)

	opts.OnConnect = func(c MQTT.Client) {
		fmt.Println("[INFO] Connected to broker.")
		if token := c.Subscribe(subTopic, 0, messageHandler); token.Wait() && token.Error() != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Subscribe failed: %v\n", token.Error())
		} else {
			fmt.Printf("[INFO] Subscribed to topic: %s\n", subTopic)
		}
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Initial connect failed: %v\n", token.Error())
		os.Exit(1)
	}

	<-ctx.Done()
	client.Disconnect(250)
	fmt.Println("[INFO] Disconnected cleanly.")
}
