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
    CheckOrigin: func(r *http.Request) bool {
        log.Printf("Origin: %s", r.Header.Get("Origin"));
        return true;
    },
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("WebSocket upgrade error: %v", err)
            return
        }
        log.Printf("WebSocket connection established from %s", r.RemoteAddr)
        defer conn.Close()

        if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"status": "Connected to server"}`)); err != nil {
            log.Printf("Initial write error: %v", err)
            return
        }

        for {
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                log.Printf("Read error: %v", err)
                break
            }
            log.Printf("Received message type %d: %s", msgType, msg)
            if msgType == websocket.TextMessage {
                if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                    log.Printf("Write error: %v", err)
                    break
                }
            }
        }
        log.Printf("WebSocket connection closed")
    })

    listenAddr := fmt.Sprintf("0.0.0.0:%s", port)
    fmt.Printf("Starting WebSocket server on %s\n", listenAddr)
    log.Fatal(http.ListenAndServe(listenAddr, nil))
}