package main

import (
	"testing"

	"github.com/EricNeid/go-relayserver/internal/test"
)

func TestNormalizePort_shouldAppend(t *testing.T) {
	// action
	result := normalizePort("8000")
	// verify
	test.Equals(t, ":8000", result)
}

func TestNormalizePort_shouldChangeNothing(t *testing.T) {
	// action
	result := normalizePort(":9000")
	// verify
	test.Equals(t, ":9000", result)
}
