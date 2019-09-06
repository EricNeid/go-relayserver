package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePort_shouldAppend(t *testing.T) {
	// action
	result := normalizePort("8000")
	// verify
	assert.Equal(t, ":8000", result)
}

func TestNormalizePort_shouldChangeNothing(t *testing.T) {
	// action
	result := normalizePort(":9000")
	// verify
	assert.Equal(t, ":9000", result)
}

func TestRunRelayServer(t *testing.T) {
	// arrange
	RunRelayServer(":8995", ":8996", "test")

	con, err := connectClient(":8996")
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

	// action
	start := time.Now()
	err = startSendingSampleStream(":8995")
	timeTrack(t, start, "TestRunRelayServer: Sending data")

	// verify
	assert.NoError(t, err)
}

func timeTrack(t *testing.T, start time.Time, name string) {
	elapsed := time.Since(start)
	t.Logf("%s took %s", name, elapsed)
}
