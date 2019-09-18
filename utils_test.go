package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRecordStream(t *testing.T) {
	// arrange
	os.Remove("testdata/recorded-sample.txt")
	unit := newStreamServer(":8080", "test")
	unit.routes()
	go func() {
		unit.listenAndServe()
	}()

	// action
	go func() {
		streamRecorded := recordStream(unit.inputStream, "testdata", "recorded-sample.txt")
		for {
			<-streamRecorded
		}
	}()
	err := sendData(":8080", "test", "test-stream")

	// verify
	ok(t, err)
	recorded, err := ioutil.ReadFile("testdata/recorded-sample.txt")
	ok(t, err)
	equals(t, "test-stream", string(recorded))

	// cleanup
	unit.shutdown()
}
