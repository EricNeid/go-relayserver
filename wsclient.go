package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	remoteAddress string
	writeStream   chan<- *[]byte
	isClosed      <-chan bool
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

func writeToConnection(conn *websocket.Conn) chan<- *[]byte {
	inputStream := make(chan *[]byte)
	go func() {
		for {
			data := <-inputStream
			if err := conn.WriteMessage(websocket.BinaryMessage, *data); err != nil {
				log.Println(err.Error())
			}
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
				isClosed <- true
				conn.Close()
				break
			}
		}
	}()
	return isClosed
}
