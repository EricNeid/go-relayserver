package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

type conn struct {
	net.Conn
}

func TestRunRelayServer(t *testing.T) {
	// arrange
	//start server
	RunRelayServer(":8081", ":8082", "")
	// connect with ws client
	ws := getConnectedWSClient(t)
	defer ws.Close()

	// action
	go postFile("testdata/SampleVideo_1280x720_30mb.mp4")
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	// verify
	fmt.Println("done")
}

func getConnectedWSClient(t *testing.T) *websocket.Conn {
	url := url.URL{
		Scheme: "ws",
		Host:   "localhost:8082",
		Path:   "/",
	}
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		assert.Fail(t, "Could not create websocket client")
	}
	return c
}

func postFile(filename string) error {
	targetURL := "http://localhost:8081/"

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}
