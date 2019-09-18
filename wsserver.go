package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// wsServer waits for websocket clients to connect.
type wsServer struct {
	server           *http.Server
	router           *http.ServeMux
	upgrader         websocket.Upgrader
	connectedClients chan *wsClient
	done             bool
}

func newWebSocketServer(port string) *wsServer {
	router := http.NewServeMux()
	server := &http.Server{Addr: port, Handler: router}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  bufferSize,
		WriteBufferSize: bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &wsServer{
		server:           server,
		router:           router,
		upgrader:         upgrader,
		connectedClients: make(chan *wsClient),
		done:             false,
	}
}

func (s *wsServer) routes() {
	s.router.HandleFunc("/clients", logRequest(s.handleClientConnect))
}

func (s *wsServer) listenAndServe() {
	s.done = false
	s.server.ListenAndServe()
}

func (s *wsServer) shutdown() {
	s.done = true                       // signal handleStream to finish reading
	time.Sleep(1000 * time.Millisecond) // give handleStream time to settle
	s.server.Shutdown(context.Background())
}

func (s *wsServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading connection: %s\n", err.Error())
		return
	}
	log.Println("WS client connected: " + ws.RemoteAddr().String())
	s.connectedClients <- &wsClient{
		remoteAddress: ws.RemoteAddr().String(),
		writeStream:   writeToConnection(ws),
		isClosed:      monitorConnection(ws),
	}
}
