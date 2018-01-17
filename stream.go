package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func waitForStream(port string, secret string) <-chan *[]byte {
	log.Println("Listening for incoming stream on " + port)

	stream := make(chan *[]byte)
	go func() {
		streamReader := http.NewServeMux()
		streamReader.HandleFunc("/"+secret, func(w http.ResponseWriter, r *http.Request) {
			log.Println("Stream connected: " + r.RemoteAddr)

			input := r.Body
			defer input.Close()

			// read from input until EOF
			buffer := make([]byte, bufferSize)
			for {
				n, err := input.Read(buffer[:cap(buffer)])
				if n == 0 {
					if err == nil {
						continue
					}
					if err == io.EOF {
						break
					}
					fmt.Println(err.Error())
				}
				chunk := buffer[:n]
				stream <- &chunk
			}
			log.Println("Stream closed")
		})
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
			// directly relay data to outstream to prevent blocking
			data := <-stream
			c <- data
			if _, err := f.Write(*data); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	return c
}
