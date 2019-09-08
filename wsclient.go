package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// wsClient represents a write-only connection to connected websocket.
type wsClient struct {
	remoteAddress string
	writeStream   chan<- *[]byte
	isClosed      <-chan bool
}

// waitForWSClients waits for connected clients. New connections are pushed on the
// returned channel.
func waitForWSClients(port string) <-chan *wsClient {
	log.Printf("Listening for incoming ws on %s/\n", port)

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
			log.Println("New ws connection: " + r.RemoteAddr)
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("WS client connected: " + ws.RemoteAddr().String())
			connectedClients <- &wsClient{
				remoteAddress: ws.RemoteAddr().String(),
				writeStream:   writeToConnection(ws),
				isClosed:      monitorConnection(ws),
			}
		})
		http.ListenAndServe(port, connectWs)
	}()

	return connectedClients
}

// writeToConnection runs a goroutine to write to the given connection.
// It returns a channel for communication.
func writeToConnection(conn *websocket.Conn) chan<- *[]byte {
	inputStream := make(chan *[]byte)
	go func() {
		for {
			data := <-inputStream
			// ignore writing error to increase throughput
			conn.WriteMessage(websocket.BinaryMessage, *data)
		}
	}()
	return inputStream
}

// monitorConnection waits for the websocket to disconnect.
// It closes the connection and send a signal on the channel.
func monitorConnection(conn *websocket.Conn) <-chan bool {
	isClosed := make(chan bool)
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Println("WS client disconnected: " + conn.RemoteAddr().String())
				isClosed <- true
				conn.Close()
				break
			}
		}
	}()
	return isClosed
}
