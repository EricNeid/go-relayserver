package relay

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WsServer waits for websocket clients to connect.
type WsServer struct {
	server          *http.Server
	router          *http.ServeMux
	upgrader        websocket.Upgrader
	IncomingClients chan *WsClient
	Port            string
	done            bool
}

// NewWebSocketServer creates new server to await websocket connections.
func NewWebSocketServer(port string) *WsServer {
	router := http.NewServeMux()
	server := &http.Server{Addr: port, Handler: router}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  bufferSize,
		WriteBufferSize: bufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WsServer{
		server:          server,
		router:          router,
		upgrader:        upgrader,
		IncomingClients: make(chan *WsClient),
		done:            false,
		Port:            port,
	}
}

// Routes configures the routes for the given server.
func (s *WsServer) Routes() {
	log.Printf("Start receiving streams on: %s/clients\n", s.Port)
	s.router.HandleFunc("/clients", logRequest(s.handleClientConnect))
}

// ListenAndServe starts listening for new websocket connections.
func (s *WsServer) ListenAndServe() {
	s.done = false
	s.server.ListenAndServe()
}

// Shutdown stops the server.
func (s *WsServer) Shutdown() {
	s.done = true                       // signal handleStream to finish reading
	time.Sleep(1000 * time.Millisecond) // give handleStream time to settle
	s.server.Shutdown(context.Background())
}

func (s *WsServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading connection: %s\n", err.Error())
		return
	}
	log.Println("WS client connected: " + ws.RemoteAddr().String())
	s.IncomingClients <- &WsClient{
		remoteAddress: ws.RemoteAddr().String(),
		writeStream:   writeToConnection(ws),
		isClosed:      monitorConnection(ws),
	}
}
