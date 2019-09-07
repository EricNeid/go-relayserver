package main

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sendData(port string, secret string, data string) error {
	url := "http://localhost" + port + "/" + secret
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func TestWaitForStream(t *testing.T) {
	// arrange
	stream := waitForStream(":8990", "test")
	go func() {
		err := sendData(":8990", "test", "Hallo, Welt")
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
	os.Remove("testdata/recorded-sample.txt")
	stream := waitForStream(":8991", "test")
	go func() {
		streamRecorded := recordStream(stream, "testdata/recorded-sample.txt")
		for {
			<-streamRecorded
		}
	}()

	// action
	err := sendData(":8991", "test", "Hallo, Welt")

	// verify
	assert.NoError(t, err)
	recorded, _ := os.Stat("testdata/recorded-sample.txt")
	assert.NotEmpty(t, recorded.Size())
}
