package relay

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/EricNeid/go-relayserver/internal/test"
)

func TestRecordStream(t *testing.T) {
	// arrange
	os.Remove("testdata/recorded-sample.txt")
	unit := NewStreamServer(":8080", "test")
	unit.Routes()
	go func() {
		unit.ListenAndServe()
	}()

	// action
	go func() {
		streamRecorded := RecordStream(unit.InputStream, "testdata", "recorded-sample.txt")
		for {
			<-streamRecorded
		}
	}()
	err := test.SendData(":8080", "test", "test-stream")

	// verify
	test.Ok(t, err)
	recorded, err := ioutil.ReadFile("testdata/recorded-sample.txt")
	test.Ok(t, err)
	test.Equals(t, "test-stream", string(recorded))

	// cleanup
	unit.Shutdown()
}
