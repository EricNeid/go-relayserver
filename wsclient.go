package relay

import (
	"log"

	"github.com/gorilla/websocket"
)

// WsClient represents a write-only connection to connected websocket.
type WsClient struct {
	remoteAddress string
	writeStream   chan<- *[]byte
	isClosed      <-chan bool
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
