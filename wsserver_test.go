package relay

import (
	"testing"

	"github.com/EricNeid/go-relayserver/internal/test"
)

func TestHandleClientConnect(t *testing.T) {
	// arrange
	unit := NewWebSocketServer(":8080")
	unit.Routes()
	go func() {
		unit.ListenAndServe()
	}()

	con, err := test.ConnectClient(":8080")
	test.Ok(t, err)
	defer con.Close()

	// action
	firstClient := <-unit.IncomingClients

	// verify
	test.Assert(t, firstClient != nil, "Connected client is nil")

	//clean
	unit.Shutdown()
}
