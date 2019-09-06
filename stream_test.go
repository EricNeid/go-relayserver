package main

import "testing"

func TestWaitForStream(t *testing.T) {
	// arrange
	var stream = waitForStream(":8989", "test")

	// action
	<-stream
}
