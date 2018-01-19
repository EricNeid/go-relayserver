package main

import (
	"net/url"
	"os/exec"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

func TestRunRelayServer(t *testing.T) {
	// arrange
	//start server
	RunRelayServer(":8081", ":8082", "")
	// connect with ws client
	ws := getConnectedWSClient(t)
	defer ws.Close()
	// channel to signal if stream was received
	done := make(chan bool)

	// action
	// start receiving data
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				break
			}
		}
		done <- true
	}()
	// start sending data
	startVideoStream(t)

	// verify
	// wait till data was send
	<-done
	// TODO compare data
}

func getConnectedWSClient(t *testing.T) *websocket.Conn {
	url := url.URL{
		Scheme: "ws",
		Host:   "localhost:8082",
		Path:   "/",
	}
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		assert.Fail(t, "Could not create websocket client: "+err.Error())
	}
	return c
}

func startVideoStream(t *testing.T) {
	c := exec.Command("testdata/stream_video.bat")
	if err := c.Run(); err != nil {
		assert.Fail(t, "Could not send video stream: "+err.Error())
	}
}
