package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func logRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		fmt.Printf("%s - streamServer - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		fn(w, r)
	}
}

// recordStream write the given stream to file. It returns the stream for further uses.
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
