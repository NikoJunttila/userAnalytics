package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}

	clientCount = 0
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	// Increment client count on connection
	clientCount++
	// Send current count to the new client
	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", clientCount)))
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// Decrement client count on disconnection
			clientCount--
			break
		}
	}
}
func handleSocketCount(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, clientCount)
}

// func main() {
// 	http.HandleFunc("/ws", handleConnections)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
