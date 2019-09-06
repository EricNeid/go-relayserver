package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePort_shouldAppend(t *testing.T) {
	// action
	result := normalizePort("8000")
	// verify
	assert.Equal(t, ":8000", result)
}

func TestNormalizePort_shouldChangeNothing(t *testing.T) {
	// action
	result := normalizePort(":9000")
	// verify
	assert.Equal(t, ":9000", result)
}
