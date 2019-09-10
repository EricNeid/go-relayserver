package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func waitForStream(port string, secret string, done <-chan bool) (*http.Server, <-chan *[]byte) {
	log.Printf("Listening for incoming stream on %s/%s\n", port, secret)

	stream := make(chan *[]byte)
	streamHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Stream client connected: " + r.RemoteAddr)

		input := r.Body
		defer input.Close()

		// read from input until EOF
		buffer := make([]byte, bufferSize)
		for {
			var readCount int
			var err error

			select {
			case <-done:
				log.Println("Stop waiting for input stream")
				return
			default:
				readCount, err = input.Read(buffer[:cap(buffer)])
			}

			if readCount > 0 {
				chunk := buffer[:readCount]
				stream <- &chunk
			}

			if err == io.EOF {
				log.Println("Stream EOF reached")
				break
			} else {
				log.Printf("Error while reading from stream: %s\n", err.Error())
				break
			}
		}
	}
	streamHandler := http.NewServeMux()
	streamHandler.HandleFunc("/"+secret, streamHandleFunc)
	streamServer := &http.Server{Addr: port, Handler: streamHandler}
	go func() {
		streamServer.ListenAndServe()
	}()

	return streamServer, stream
}

// recordStream write the given stream to file. It returns the stream for further uses
// and is not blocking the channel.
func recordStream(stream <-chan *[]byte, path string, file string) <-chan *[]byte {
	c := make(chan *[]byte)
	os.MkdirAll(path, os.ModePerm)
	f, err := os.Create(filepath.Join(path, file))
	if err != nil {
		log.Println(err.Error())
		return stream
	}

	go func() {
		defer f.Close()
		for {
			newChunk := <-stream
			c <- newChunk
			if _, err := f.Write(*newChunk); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	return c
}
