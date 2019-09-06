package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelayStreamToWSClients(t *testing.T) {
	// arrange
	stream := waitForStream(":8993", "test")
	clients := waitForWSClients(":8994")
	con, err := connectClient(":8994")
	if err != nil {
		assert.Fail(t, "Error while creating client connection")
	}
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
	err = startSendingSampleStream(":8993")
	timeTrack(t, start, "TestRelayStreamToWSClients: Sending data")

	// verify
	assert.NoError(t, err)
}
