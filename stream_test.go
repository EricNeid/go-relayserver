package main

import (
	"os"
	"testing"
)

func TestWaitForStream(t *testing.T) {
	// arrange
	stream := waitForStream(":8990", "test")
	go func() {
		err := sendData(":8990", "test", "Hallo, Welt")
		ok(t, err)
	}()

	// action
	firstChunk := <-stream

	// verify
	equals(t, true, firstChunk != nil)
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
	ok(t, err)
	recorded, _ := os.Stat("testdata/recorded-sample.txt")
	equals(t, true, recorded.Size() > 0)
}
