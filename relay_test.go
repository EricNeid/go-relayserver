package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelayStreamToWSClients(t *testing.T) {
	// arrange
	stream := waitForStream(":8989", "test")
	clients := waitForWSClients(":8990")
	con, err := connectClient()
	if err != nil {
		assert.Fail(t, "Error while creating client connection")
	}
	defer con.Close()

	go func() {
		for {
			con.ReadMessage()
			//fmt.Printf("Received %d bytes", len(bytes))
			//if errClient != nil {
			//	assert.Fail(t, "Error while receiving bytes %s", errClient)
			//}
		}
	}()

	relayStreamToWSClients(stream, clients)

	// action
	err = startSendingSampleStream(":8989")

	// verify
	assert.NoError(t, err)
}
