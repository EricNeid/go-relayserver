package test

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// Assert fails the test if the condition is false.
func Assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

// SendData sends test data to localhost with given port and secret.
func SendData(port string, secret string, data string) error {
	url := "http://localhost" + port + "/stream/" + secret
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// SendVideo send test video to localhost with given port.
func SendVideo(port string) error {
	streamSender := exec.Command("ffmpeg",
		"-i", "testdata/sample.mp4",
		"-f", "mpegts",
		"-codec:v", "mpeg1video",
		"-s", "1280x720",
		"-rtbufsize", "2048M",
		"-r", "30",
		"-b:v", "3000k",
		"-q:v", "6",
		"http://localhost"+port+"/stream/test")

	_, err := streamSender.Output()
	return err
}

// TimeTrack prints the elapsed time since start in seconds, together with the given name.
func TimeTrack(t *testing.T, start time.Time, name string) {
	elapsed := time.Since(start)
	t.Logf("%s took %s", name, elapsed)
}

// ConnectClient connects a websocket to local test server on the given port.
func ConnectClient(port string) (*websocket.Conn, error) {
	url := "ws://localhost" + port + "/clients"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	return c, err
}
