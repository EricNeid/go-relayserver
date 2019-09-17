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

/*
func TestHandleStream_interruptStream(t *testing.T) {
	// arrange
	largeFile, _ := ioutil.ReadFile("testdata/sample-data.txt")
	unit := newStreamServer(":8080", "test")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()
	go func() {
		err := sendData(":8080", "test", string(largeFile)+string(largeFile)+string(largeFile))
		ok(t, err)
	}()

	// action
	received := string(*<-unit.inputStream)
	unit.shutdown()

	// verify
	equals(t, "Hallo, Welt", received)
}
*/
/*
func TestWaitForStream(t *testing.T) {
	// arrange
	done := make(chan bool, 1)
	server, stream := waitForStream(":8990", "test", done)

	// action
	go func() {
		err := sendData(":8990", "test", "Hallo, Welt")
		ok(t, err)
	}()

	// verify
	received := <-stream
	assert(t, received != nil, "Received data is null")
	equals(t, "Hallo, Welt", string(*received))

	// action
	go func() {
		err := sendData(":8990", "test", "Hallo, Welt 2")
		ok(t, err)
	}()

	// verify
	received = <-stream
	assert(t, received != nil, "Received data is null")
	equals(t, "Hallo, Welt 2", string(*received))

	// cleanup
	done <- true
	server.Shutdown(context.Background())
}

func TestRecordStream(t *testing.T) {
	// arrange
	done := make(chan bool, 1)
	os.Remove("testdata/recorded-sample.txt")
	server, stream := waitForStream(":8990", "test", done)
	go func() {
		streamRecorded := recordStream(stream, "testdata", "recorded-sample.txt")
		for {
			<-streamRecorded
		}
	}()

	// action
	err := sendData(":8990", "test", "Hallo, Welt")

	// verify
	ok(t, err)
	recorded, err := ioutil.ReadFile("testdata/recorded-sample.txt")
	ok(t, err)
	equals(t, "Hallo, Welt", string(recorded))

	// cleanup
	server.Shutdown(context.Background())
	done <- true
}

*/
