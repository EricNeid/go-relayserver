package main

import (
	"testing"
)

func TestHandleClientConnect(t *testing.T) {
	// arrange
	unit := newWebSocketServer(":8080")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()

	con, err := connectClient(":8080")
	ok(t, err)
	defer con.Close()

	// action
	firstClient := <-unit.connectedClients

	// verify
	assert(t, firstClient != nil, "Connected client is nil")

	//clean
	unit.shutdown()
}
