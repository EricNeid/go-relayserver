package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func startSendingSampleStream() error {
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

	_, err := streamSender.Output()
	return err
}

func TestWaitForStream(t *testing.T) {
	// arrange
	stream := waitForStream(":8989", "test")
	go func() {
		err := startSendingSampleStream()
		if err != nil {
			assert.Fail(t, "Error while sending stream")
		}
	}()

	// action
	firstChunk := <-stream

	// verify
	assert.NotEmpty(t, firstChunk)
}

func TestRecordStream(t *testing.T) {
	// arrange
	os.Remove("testdata/recorded-sample.mpeg")
	stream := waitForStream(":8989", "test")
	go func() {
		streamRecorded := recordStream(stream, "testdata/recorded-sample.mpeg")
		for {
			<-streamRecorded
		}
	}()

	// action
	err := startSendingSampleStream()

	// verify
	assert.NoError(t, err)
}
