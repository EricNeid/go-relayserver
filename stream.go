package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func waitForStream(port string, secret string, done <-chan bool) <-chan *[]byte {
	log.Printf("Listening for incoming stream on %s/%s\n", port, secret)

	stream := make(chan *[]byte)
	streamReader := http.NewServeMux()
	streamReader.HandleFunc("/"+secret, func(w http.ResponseWriter, r *http.Request) {
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
	})
	go func() {
		http.ListenAndServe(port, streamReader)
	}()
	return stream
}

// recordStream write the given stream to file. It returns the stream for further uses
// and is not blocking the channel.
func recordStream(stream <-chan *[]byte, recordName string) <-chan *[]byte {
	c := make(chan *[]byte)
	f, err := os.Create(recordName)
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
