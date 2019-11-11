package relay

import (
	"testing"
	"time"

	"github.com/EricNeid/go-relayserver/internal/test"
)

func TestStreamToWSClients(t *testing.T) {
	// arrange
	streamServer := NewStreamServer(":8080", "test")
	streamServer.Routes()

	webSocketServer := NewWebSocketServer(":8081")
	webSocketServer.Routes()

	go func() {
		streamServer.ListenAndServe()
	}()
	go func() {
		webSocketServer.ListenAndServe()
	}()

	con, err := test.ConnectClient(":8081")
	test.Ok(t, err)
	defer con.Close()
	go func() {
		for {
			_, _, err := con.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	// action
	StreamToWSClients(streamServer.InputStream, webSocketServer.IncomingClients)
	start := time.Now()
	err = test.SendData(":8080", "test", "Hallo, Welt")
	test.TimeTrack(t, start, "TestRelayStreamToWSClients: Sending data")

	// verify
	test.Ok(t, err)

	// cleanup
	streamServer.Shutdown()
	webSocketServer.Shutdown()
}
