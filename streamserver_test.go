package relay

import (
	"testing"

	"github.com/EricNeid/go-relayserver/internal/test"
)

func TestHandleStream(t *testing.T) {
	// arrange
	unit := NewStreamServer(":8080", "test")
	unit.Routes()
	go func() {
		unit.ListenAndServe()
	}()

	// action
	go func() {
		err := test.SendData(":8080", "test", "Hallo, Welt")
		test.Ok(t, err)
	}()

	// verify
	received := string(*<-unit.InputStream)
	test.Equals(t, "Hallo, Welt", received)

	// cleanup
	unit.Shutdown()
}

func TestHandleStream_twoStreamsSequential(t *testing.T) {
	// arrange
	unit := NewStreamServer(":8080", "test")
	unit.Routes()
	go func() {
		unit.ListenAndServe()
	}()

	// action
	go func() {
		err := test.SendData(":8080", "test", "Hallo, Welt")
		test.Ok(t, err)
	}()

	// verify
	received := string(*<-unit.InputStream)
	test.Equals(t, "Hallo, Welt", received)

	// action 2
	go func() {
		err := test.SendData(":8080", "test", "Was gibt's?")
		test.Ok(t, err)
	}()

	// verify 2
	received = string(*<-unit.InputStream)
	test.Equals(t, "Was gibt's?", received)

	// cleanup
	unit.Shutdown()
}

func TestHandleStream_twoStreamsParallel(t *testing.T) {
	// arrange
	unit := NewStreamServer(":8080", "test")
	unit.Routes()
	go func() {
		unit.ListenAndServe()
	}()

	// action
	go func() {
		err := test.SendData(":8080", "test", "test-stream-1")
		test.Ok(t, err)
	}()
	go func() {
		err := test.SendData(":8080", "test", "test-stream-2")
		test.Ok(t, err)
	}()

	// verify
	received := string(*<-unit.InputStream)
	test.Equals(t, "test-stream-1", received)
	received2 := string(*<-unit.InputStream)
	test.Equals(t, "test-stream-2", received2)

	// cleanup
	unit.Shutdown()
}

func TestHandleStream_interruptStream(t *testing.T) {
	// arrange
	unit := NewStreamServer(":8080", "test")
	unit.Routes()
	go func() {
		unit.ListenAndServe()
	}()
	go func() {
		test.SendVideo(":8080")
	}()

	// action
	received := *<-unit.InputStream
	go func() {
		<-unit.InputStream
	}()
	unit.Shutdown() // input reading from stream

	// verify
	test.Assert(t, len(received) > 0, "Received chunk is empty")
}
