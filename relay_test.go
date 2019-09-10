package main

import (
	"testing"
	"time"
)

func TestRelayStreamToWSClients(t *testing.T) {
	// arrange
	done := make(chan bool, 1)
	_, stream := waitForStream(":8993", "test", done)
	clients := waitForWSClients(":8994")
	con, err := connectClient(":8994")
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

	relayStreamToWSClients(stream, clients)

	// action
	start := time.Now()
	err = sendData(":8993", "test", "Hallo, Welt")
	timeTrack(t, start, "TestRelayStreamToWSClients: Sending data")

	// verify
	ok(t, err)

	//cleanup
	done <- true
}
