package main

import (
    "encoding/json"
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
                var data map[string]interface{}
                if err := json.Unmarshal(msg, &data); err != nil {
                    log.Printf("Unmarshal error: %v", err)
                    continue
                }

                if data["type"] == "file" {
                    // Step 1: Transfer wallpaper data
                    transferMsg := map[string]interface{}{
                        "type":     "transfer",
                        "endpoint": 1,
                        "data":     data["data"],
                    }
                    transferJSON, _ := json.Marshal(transferMsg)
                    if err := conn.WriteMessage(websocket.TextMessage, transferJSON); err != nil {
                        log.Printf("Write error: %v", err)
                        break
                    }

                    // Step 2: Control command to write file (SparseRestore placeholder)
                    controlMsg1 := map[string]interface{}{
                        "type":        "control",
                        "requestType": "vendor",
                        "recipient":   "device",
                        "request":     0x40, // SparseRestore filesystem write (needs exact value)
                        "value":       0x01, // Write operation
                        "index":       0,
                    }
                    controlJSON1, _ := json.Marshal(controlMsg1)
                    if err := conn.WriteMessage(websocket.TextMessage, controlJSON1); err != nil {
                        log.Printf("Write error: %v", err)
                        break
                    }

                    // Step 3: Control command to refresh SpringBoard (SparseRestore placeholder)
                    controlMsg2 := map[string]interface{}{
                        "type":        "control",
                        "requestType": "vendor",
                        "recipient":   "device",
                        "request":     0x40, // SparseRestore refresh (needs exact value)
                        "value":       0x02, // Refresh operation
                        "index":       0,
                    }
                    controlJSON2, _ := json.Marshal(controlMsg2)
                    if err := conn.WriteMessage(websocket.TextMessage, controlJSON2); err != nil {
                        log.Printf("Write error: %v", err)
                        break
                    }

                    // Step 4: Notify completion
                    if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"status": "Wallpaper applied"}`)); err != nil {
                        log.Printf("Write error: %v", err)
                        break
                    }
                }
            }
        }
        log.Printf("WebSocket connection closed")
    })

    listenAddr := fmt.Sprintf("0.0.0.0:%s", port)
    fmt.Printf("Starting WebSocket server on %s\n", listenAddr)
    log.Fatal(http.ListenAndServe(listenAddr, nil))
}