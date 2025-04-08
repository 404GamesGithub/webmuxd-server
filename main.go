package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// Use Render's PORT environment variable, fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define WebSocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// Handle incoming messages (e.g., .tendies file)
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				break
			}

			// Process message (simplified for wallpaper)
			if msgType == websocket.TextMessage {
				log.Printf("Received: %s", msg)
				// Example: Send back a command to write the file via WebUSB
				response := []byte(`{"type": "transfer", "endpoint": 1, "data": ` + string(msg) + `}`)
				err = conn.WriteMessage(websocket.TextMessage, response)
				if err != nil {
					log.Printf("Write error: %v", err)
					break
				}
			}
		}
	})

	// Bind to 0.0.0.0 and PORT
	listenAddr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Starting WebSocket server on %s\n", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}