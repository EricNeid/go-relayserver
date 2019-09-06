package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitForStream(t *testing.T) {
	// arrange
	streamReceiver := waitForStream(":8989", "test")
	streamSender := exec.Command("ffmpeg",
		"-i", "testdata/sample.mp4",
		"-f", "mpegts",
		"-codec:v", "mpeg1video",
		"-s", "1280x720",
		"-rtbufsize", "2048M",
		"-r", "30",
		"-b:v", "3000k",
		"-q:v", "6",
		"http://localhost:8989/test")

	go func() {
		_, err := streamSender.Output()
		if err != nil {
			assert.Fail(t, "Error while sending stream")
		}
	}()

	// action
	firstChunk := <-streamReceiver

	// verify
	assert.NotEmpty(t, firstChunk)
}
