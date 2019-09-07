package main

import (
	"testing"
)

func TestWaitForWSClients(t *testing.T) {
	// arrange
	clients := waitForWSClients(":8992")

	con, err := connectClient(":8992")
	ok(t, err)
	defer con.Close()

	// action
	firstClient := <-clients

	// verify
	assert(t, firstClient != nil, "Connected client is nil")
}
