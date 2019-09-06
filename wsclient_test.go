package main

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func connectClient() (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8990", nil)
	return c, err
}

func TestWaitForWSClients(t *testing.T) {
	// arrange
	clients := waitForWSClients(":8990")

	con, err := connectClient()
	if err != nil {
		assert.Fail(t, "Error while creating client connection")
	}
	defer con.Close()

	// action
	firstClient := <-clients

	// verify
	assert.NotNil(t, firstClient)
}
