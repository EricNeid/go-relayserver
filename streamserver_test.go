package main

import (
	"testing"
)

func TestHandleStream(t *testing.T) {
	// arrange
	unit := newStreamServer(":8080", "test")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()

	// action
	go func() {
		err := sendData(":8080", "test", "Hallo, Welt")
		ok(t, err)
	}()

	// verify
	received := string(*<-unit.inputStream)
	equals(t, "Hallo, Welt", received)

	// cleanup
	unit.shutdown()
}

func TestHandleStream_twoStreamsSequential(t *testing.T) {
	// arrange
	unit := newStreamServer(":8080", "test")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()

	// action
	go func() {
		err := sendData(":8080", "test", "Hallo, Welt")
		ok(t, err)
	}()

	// verify
	received := string(*<-unit.inputStream)
	equals(t, "Hallo, Welt", received)

	// action 2
	go func() {
		err := sendData(":8080", "test", "Was gibt's?")
		ok(t, err)
	}()

	// verify 2
	received = string(*<-unit.inputStream)
	equals(t, "Was gibt's?", received)

	// cleanup
	unit.shutdown()
}

func TestHandleStream_twoStreamsParallel(t *testing.T) {
	// arrange
	unit := newStreamServer(":8080", "test")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()

	// action
	go func() {
		err := sendData(":8080", "test", "test-stream-1")
		ok(t, err)
	}()
	go func() {
		err := sendData(":8080", "test", "test-stream-2")
		ok(t, err)
	}()

	// verify
	received := string(*<-unit.inputStream)
	equals(t, "test-stream-1", received)
	received2 := string(*<-unit.inputStream)
	equals(t, "test-stream-2", received2)

	// cleanup
	unit.shutdown()
}

func TestHandleStream_interruptStream(t *testing.T) {
	// arrange
	unit := newStreamServer(":8080", "test")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()
	go func() {
		sendVideo(":8080")
	}()

	// action
	received := *<-unit.inputStream
	go func() {
		<-unit.inputStream
	}()
	unit.shutdown() // input reading from stream

	// verify
	assert(t, len(received) > 0, "Received chunk is empty")
}
