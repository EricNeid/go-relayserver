package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// assert fails the test if the condition is false.
func assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

// sendData sends test data to localhost with given port and secret.
func sendData(port string, secret string, data string) error {
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

// timeTrack prints the elapsed time since start in seconds, together with the given name.
func timeTrack(t *testing.T, start time.Time, name string) {
	elapsed := time.Since(start)
	t.Logf("%s took %s", name, elapsed)
}

func connectClient(port string) (*websocket.Conn, error) {
	url := "ws://localhost" + port
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	return c, err
}
