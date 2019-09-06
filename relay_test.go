package main

import (
	"testing"

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
	err = startSendingSampleStream(":8993")

	// verify
	assert.NoError(t, err)
}

