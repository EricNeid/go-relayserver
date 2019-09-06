package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func startSendingSampleStream(port string) error {
	streamSender := exec.Command("ffmpeg",
		"-i", "testdata/sample.mp4",
		"-f", "mpegts",
		"-codec:v", "mpeg1video",
		"-s", "1280x720",
		"-rtbufsize", "2048M",
		"-r", "30",
		"-b:v", "3000k",
		"-q:v", "6",
		"http://localhost"+port+"/test")

	_, err := streamSender.Output()
	return err
}

func TestWaitForStream(t *testing.T) {
	// arrange
	stream := waitForStream(":8990", "test")
	go func() {
		err := startSendingSampleStream(":8990")
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
	stream := waitForStream(":8991", "test")
	go func() {
		streamRecorded := recordStream(stream, "testdata/recorded-sample.mpeg")
		for {
			<-streamRecorded
		}
	}()

	// action
	err := startSendingSampleStream(":8991")

	// verify
	assert.NoError(t, err)
	recorded, _ := os.Stat("testdata/recorded-sample.mpeg")
	assert.NotEmpty(t, recorded.Size())
}
