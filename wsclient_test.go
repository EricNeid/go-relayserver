package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func connectClient(port string) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost"+port, nil)
	return c, err
}

func TestWaitForWSClients(t *testing.T) {
	// arrange
	clients := waitForWSClients(":8992")

	con, err := connectClient(":8992")
	ok(t, err)
	defer con.Close()

	// action
	firstClient := <-clients

	// verify
	equals(t, true, firstClient != nil)
}
