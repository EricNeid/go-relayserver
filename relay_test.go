package main

import (
	"testing"
	"time"
)

func TestRelayStreamToWSClients(t *testing.T) {
	// arrange
	streamServer := newStreamServer(":8080", "test")
	streamServer.routes()

	webSocketServer := newWebSocketServer(":8081")
	webSocketServer.routes()

	go func() {
		streamServer.listenAndServe()
	}()
	go func() {
		webSocketServer.listenAndServe()
	}()

	con, err := connectClient(":8081")
	ok(t, err)
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
	relayStreamToWSClients(streamServer.inputStream, webSocketServer.incomingClients)
	start := time.Now()
	err = sendData(":8080", "test", "Hallo, Welt")
	timeTrack(t, start, "TestRelayStreamToWSClients: Sending data")

	// verify
	ok(t, err)

	// cleanup
	streamServer.shutdown()
	webSocketServer.shutdown()
}
