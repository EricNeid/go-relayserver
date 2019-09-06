package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelayStreamToWSClients(t *testing.T) {
	// arrange
	stream := waitForStream(":8990", "test")
	clients := waitForWSClients(":8991")
	con, err := connectClient(":8991")
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
	err = startSendingSampleStream(":8990")

	// verify
	assert.NoError(t, err)
}
