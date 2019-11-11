package relay

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

const timeFormat string = "Mon Jan 2 15:04:05 2006"

// StreamServer represents a server, ready to access a single input stream.
// To create a new instance use s := newStreamServer(":8080", "MySecret")
// Before using the server, you need to call s.routes() to configure the routes.
// To shutdown the server use s.shutdown()
type StreamServer struct {
	server            *http.Server
	router            *http.ServeMux
	isStreamConnected chan bool
	InputStream       chan *[]byte
	done              bool
	Secret            string
	Port              string
}

// NewStreamServer creates new instance of stream server.
// It is regsitered on the given port and access to stream is protected by
// given secret.
func NewStreamServer(port string, secret string) *StreamServer {
	router := http.NewServeMux()
	server := &http.Server{Addr: port, Handler: router}

	return &StreamServer{
		server:            server,
		router:            router,
		isStreamConnected: make(chan bool, 1),
		InputStream:       make(chan *[]byte),
		Secret:            secret,
		Port:              port,
		done:              false,
	}
}

// Routes registers function handler for this SteamServer.
func (s *StreamServer) Routes() {
	log.Printf("Start receiving streams on: %s/stream/%s\n", s.Port, s.Secret)
	s.router.HandleFunc("/stream/"+s.Secret, logRequest(s.handleStream))
}

// ListenAndServe starts listening for new stream.
func (s *StreamServer) ListenAndServe() {
	s.done = false
	s.server.ListenAndServe()
}

// Shutdown stops server.
func (s *StreamServer) Shutdown() {
	s.done = true                       // signal handleStream to finish reading
	time.Sleep(1000 * time.Millisecond) // give handleStream time to settle
	s.server.Shutdown(context.Background())
}

func (s *StreamServer) handleStream(w http.ResponseWriter, r *http.Request) {
	s.isStreamConnected <- true

	input := r.Body
	defer input.Close()

	buffer := make([]byte, bufferSize)
	for !s.done {
		var readCount int
		var err error

		readCount, err = input.Read(buffer[:cap(buffer)])

		if readCount > 0 {
			chunk := buffer[:readCount]
			s.InputStream <- &chunk
		}

		if err == io.EOF {
			log.Println("Stream EOF reached")
			break
		} else if err != nil {
			log.Printf("Error while reading from stream: %s\n", err.Error())
			break
		}
	}
	log.Println("Stop waiting for input stream")
	<-s.isStreamConnected
}
