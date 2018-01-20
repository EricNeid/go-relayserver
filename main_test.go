package main

import (
	"io/ioutil"
	"net/url"
	"os/exec"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

func TestRunRelayServer__ShouldReceiveWholeStream(t *testing.T) {
	// arrange
	//start server
	RunRelayServer(":8081", ":8082", "secret1234")
	// connect with ws client
	ws := getConnectedWSClient(t)
	// channel to signal if stream was received
	done := make(chan bool)
	receivedBytes := make(chan []byte)

	// action
	// start receiving data
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				break
			}
			//receivedBytes <- chunk
		}
	}()
	// start sending data
	startVideoStream(t, done)

	// verify
	// wait till data was send
	<-done
	ws.Close()
	close(receivedBytes)

	expected := getBytesForStream(t)
	received := joinChannel(receivedBytes)
	assert.Equal(t, len(expected), len(received))
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

func getBytesForStream(t *testing.T) []byte {
	content, err := ioutil.ReadFile("testdata/SampleVideo_1280x720_30mb.mp4")
	if err != nil {
		assert.FailNow(t, "Could not read video stream: "+err.Error())
	}
	return content
}

func joinChannel(chunks <-chan []byte) []byte {
	var joined []byte
	for chunk := range chunks {
		for _, b := range chunk {
			joined = append(joined, b)
		}
	}
	return joined
}
