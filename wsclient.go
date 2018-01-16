package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	writeStream chan<- []byte
	isClosed    <-chan bool
}

func waitForWSClients(port string) <-chan *wsClient {
	connectedClients := make(chan *wsClient)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  bufferSize,
		WriteBufferSize: bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	go func() {
		connectWs := http.NewServeMux()
		connectWs.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Received ws connection from: " + r.RemoteAddr)

			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Client connected: " + ws.RemoteAddr().String())

			connectedClients <- &wsClient{}
		})
		http.ListenAndServe(port, connectWs)
	}()

	return connectedClients
}
