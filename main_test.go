package main

import (
	"testing"
	"time"
)

func TestNormalizePort_shouldAppend(t *testing.T) {
	// action
	result := normalizePort("8000")
	// verify
	equals(t, ":8000", result)
}

func TestNormalizePort_shouldChangeNothing(t *testing.T) {
	// action
	result := normalizePort(":9000")
	// verify
	equals(t, ":9000", result)
}

func TestRunRelayServer(t *testing.T) {
	// arrange
	RunRelayServer(":8995", ":8996", "test", false)

	con, err := connectClient(":8996")
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
	start := time.Now()
	err = sendData(":8995", "test", "Hallo, Welt")
	timeTrack(t, start, "TestRunRelayServer: Sending data")

	// verify
	ok(t, err)
}
