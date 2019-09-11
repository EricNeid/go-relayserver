package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const timeFormat string = "Mon Jan 2 15:04:05 2006"

type streamServer struct {
	router            *http.ServeMux
	isStreamConnected chan bool
	inputStream       chan *[]byte
	done              chan bool
	secret            string
}

func newStreamServer(secret string) *streamServer {
	router := http.NewServeMux()

	return &streamServer{
		router:            router,
		isStreamConnected: make(chan bool, 1),
		inputStream:       make(chan *[]byte),
		secret:            secret,
	}
}

func (s *streamServer) routes() {
	s.router.HandleFunc("/stream/"+s.secret, s.log(s.handleStream))
}

func (s *streamServer) handleStream(w http.ResponseWriter, r *http.Request) {
	s.isStreamConnected <- true

	input := r.Body
	defer input.Close()

	buffer := make([]byte, bufferSize)
	for {
		var readCount int
		var err error

		isDone := false

		select {
		case <-s.done:
			isDone = true
		default:
			readCount, err = input.Read(buffer[:cap(buffer)])
		}

		if isDone {
			log.Println("Stop waiting for input stream")
			break
		}

		if readCount > 0 {
			chunk := buffer[:readCount]
			s.inputStream <- &chunk
		}

		if err == io.EOF {
			log.Println("Stream EOF reached")
			break
		} else if err != nil {
			log.Printf("Error while reading from stream: %s\n", err.Error())
			break
		}
	}
	<-s.isStreamConnected
}

func (s *streamServer) log(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		fmt.Printf("%s - streamServer - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		fn(w, r)
	}
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
