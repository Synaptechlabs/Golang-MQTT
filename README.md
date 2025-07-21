# Golang-MQTT

Lightweight MQTT listener and message router written in Go. Originally built as a prototype for commercial IoT systems in RV automation, this project now serves as a minimal, modern proof-of-concept for MQTT integration between a browser frontend and a backend service.

> Part of the Arctic Code Vault — originally published as [github.com/NGC3031](https://github.com/NGC3031). Maintained by Scott Douglass, SynaptechLabs.ai.

---

## 💡 Overview

This project demonstrates:

- Real-time MQTT message listening and publishing via [Eclipse Paho MQTT](https://github.com/eclipse/paho.mqtt.golang)
- A minimal HTML+JavaScript MQTT web client using `paho-mqtt` (over WebSocket)
- Lightweight design for fast testing, debugging, and custom IoT message routing

---

## 🧩 Architecture

The backend uses a simple **Model-View-Controller (MVC)** pattern:

- **Model:** Internal message state counter and topic filtering logic
- **View:** Web frontend using a basic HTML form + Bootstrap UI
- **Controller:** MQTT client connects to broker, listens and republishes messages

---

## ⚙️ Technologies

- **Language:** Go (Golang)
- **MQTT Library:** `github.com/eclipse/paho.mqtt.golang`
- **Frontend MQTT JS:** [`paho-mqtt`](https://cdnjs.com/libraries/paho-mqtt)
- **Websocket Server:** Legacy Go `websocket` (not `gorilla` for simplicity)

---

## 🚀 Quick Start

### 1. Clone & Run the Go Server

```bash
git clone https://github.com/YOUR_NEW_REPO/Golang-MQTT.git
cd Golang-MQTT
go run mqtt.go
```

This will:
- Connect to the public HiveMQ broker (`broker.hivemq.com:1883`)
- Subscribe to all messages under topic `test4472/#`
- Re-publish responses to `test4472/result`

---

### 2. Open the Web Client

Open `sock.html` in your browser (no server needed).

- Send a message to `test4472/input`
- The Go server receives and publishes a response to `test4472/result`
- Your browser client will auto-display received messages

---

## 🔐 Broker Setup

- **Broker:** `broker.hivemq.com`
- **Port:** `1883` (Go) and `8000` (WebSocket for browser)
- **SSL:** Not required (test mode only)
- **Topic:** Use isolated topic like `test4472` to avoid public flood

---

## 🛠️ Customization

- Update topic prefixes in `mqtt.go` and `sock.html` to suit your device namespace
- Add JSON handling or routing logic in the Go `messageHandler`
- Swap in `gorilla/websocket` for more advanced WebSocket support if needed

---

## 📜 Legacy Notes

Originally published in 2019 as part of the **NGC3031** GitHub projects. The base prototype enabled real-time IoT message flow from RV sensors to a lightweight dashboard, with MQTT levels routed and visually represented in the browser.

---

## ⚠️ Disclaimer

> This is proof-of-concept code. Use in production at your own risk. No warranties, but lots of ❤️ and learning.

---

## ✍️ Author

**Scott Douglass**  
[SynaptechLabs.ai](https://synaptechlabs.ai) — building neuro-inspired, lightweight systems for edge and embedded AI.
