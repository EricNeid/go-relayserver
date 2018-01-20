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
	RunRelayServer(":8081", ":8082", "secret1234")
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
	}()
	// start sending data
	startVideoStream(t, done)

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
		assert.FailNow(t, "Could not create websocket client: "+err.Error())
	}
	return c
}

func startVideoStream(t *testing.T, done chan<- bool) {
	go func() {
		c := exec.Command("testdata\\stream_video.bat")
		err := c.Run()
		if err != nil {
			assert.FailNow(t, "Could not send video stream: "+err.Error())
		}
		done <- true
	}()
}
