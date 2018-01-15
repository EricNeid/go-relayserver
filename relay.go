package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const recordToFile = false

var recordName = fmt.Sprintf("%s_record.mpeg", time.Now().Format("20060102_1504"))

const portStream = ":8081"
const portWS = ":8082"

const bufferSize = 4 * 1000 * 1024 // 4MB

var upgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connectedWSClients []*websocket.Conn

var stream = make(chan []byte)
var done = make(chan bool)

func main() {
	go waitForWebSockets()
	go waitForStream()
	go workerThread()

	fmt.Println("Relay started, hit Enter-key to close")

	fmt.Scanln()
	done <- true

	fmt.Println("Shuting down...")
}

func waitForStream() {
	log.Println("Listening for incoming stream on " + portStream)

	serverVideoStream := http.NewServeMux()
	serverVideoStream.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received stream connection from: " + r.RemoteAddr)

		input := r.Body
		defer input.Close()

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
			stream <- chunk
		}
		log.Println("Stream closed")
	})
	http.ListenAndServe(portStream, serverVideoStream)
}

func waitForWebSockets() {
	log.Println("Listening for incoming ws on " + portWS)

	serverWS := http.NewServeMux()
	serverWS.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received ws connection from: " + r.RemoteAddr)

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Client subscribed: " + r.RemoteAddr)

		connectedWSClients = append(connectedWSClients, ws)
	})
	http.ListenAndServe(portWS, serverWS)
}

func workerThread() {
	var recording *os.File
	if recordToFile {
		f, err := os.Create(recordName)
		if err != nil {
			log.Println(err.Error())
		}
		defer f.Close()
		recording = f
	}

	for {
		select {
		case data := <-stream:
			if recordToFile {
				go writeToFile(data, recording)
			}
			for _, conn := range connectedWSClients {
				if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					log.Println(err.Error())
				}
			}
		case <-done:
			for _, c := range connectedWSClients {
				c.Close()
			}
		}
	}
}

func writeToFile(data []byte, file *os.File) {
	if file != nil {
		file.Write(data)
	}
}
