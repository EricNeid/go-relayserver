package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

const testVideo = "testdata/SampleVideo_1280x720_5mb.mp4"
const streamVideo = "testdata\\stream_video.bat"

func TestRunRelayServer__ShouldReceiveWholeStream(t *testing.T) {
	// arrange
	//start server
	RunRelayServer(":8081", ":8082", "secret1234")
	// connect with ws client
	ws := getConnectedWSClient(t)
	// channel to signal if stream was received
	done := make(chan bool)
	var receivedBytes [][]byte

	// action
	// start receiving data
	go func() {
		for {
			_, chunk, err := ws.ReadMessage()
			if err != nil {
				break
			}
			receivedBytes = append(receivedBytes, chunk)
		}
	}()
	// start sending data
	//startVideoStream(t, done)
	streamTestData(t, done)

	// verify
	// wait till data was send
	<-done
	ws.Close()

	expected := getBytesForStream(t)
	received := joinBytes(receivedBytes)
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
		c := exec.Command(streamVideo)
		err := c.Run()
		if err != nil {
			assert.FailNow(t, "Could not send video stream: "+err.Error())
		}
		done <- true
	}()
}

func streamTestData(t *testing.T, done chan<- bool) {
	input, err := os.Open(testVideo)
	if err != nil {
		assert.FailNow(t, "Could not open test video: "+err.Error())
	}

	req, err := http.NewRequest("Post", "localhost:8081/secret1234", input)
	if err != nil {
		assert.FailNow(t, "Could not create request: "+err.Error())
	}
	client := http.DefaultClient
	go func() {
		_, err := client.Do(req)
		if err != nil {
			assert.FailNow(t, "Could not send request: "+err.Error())
		}
		done <- true
	}()
}

func getBytesForStream(t *testing.T) []byte {
	content, err := ioutil.ReadFile(testVideo)
	if err != nil {
		assert.FailNow(t, "Could not read video stream: "+err.Error())
	}
	return content
}

func joinBytes(chunks [][]byte) []byte {
	var joined []byte
	for _, chunk := range chunks {
		for _, b := range chunk {
			joined = append(joined, b)
		}
	}
	return joined
}
