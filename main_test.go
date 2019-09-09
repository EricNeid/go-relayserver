package main

import (
	"testing"
)

func TestNormalizePort_shouldAppend(t *testing.T) {
	// action
	result := normalizePort("8000")
	// verify
	equals(t, ":8000", result)
}

func TestNormalizePort_shouldChangeNothing(t *testing.T) {
	// action
	result := normalizePort(":9000")
	// verify
	equals(t, ":9000", result)
}
