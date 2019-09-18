package main

import (
	"fmt"
	"net/http"
	"time"
)

func logRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		fmt.Printf("%s - streamServer - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		fn(w, r)
	}
}
