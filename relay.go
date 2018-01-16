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

const portStream = ":8081"
const portWS = ":8082"
const bufferSize = 4 * 1000 * 1024 // 4MB

var wsClients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var recordName = fmt.Sprintf("%s_record.mpeg", time.Now().Format("20060102_1504"))

func main() {
	stream := make(chan []byte)
	done := make(chan bool)

	var recording *os.File
	if recordToFile {
		f, err := os.Create(recordName)
		if err != nil {
			log.Println(err.Error())
		}
		defer f.Close()
		recording = f
	}

	go waitForStream(portStream, stream)
	go consumeStream(stream, done, recording)

	fmt.Println("Relay started, hit Enter-key to close")

	fmt.Scanln()
	done <- true

	fmt.Println("Shuting down...")
}

func waitForStream(port string, stream chan<- []byte) {
	log.Println("Listening for incoming stream on " + portStream)

	streamReader := http.NewServeMux()
	streamReader.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received stream connection from: " + r.RemoteAddr)

		input := r.Body
		defer input.Close()

		// read from input intil EOF
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
	http.ListenAndServe(port, streamReader)
}

func waitForWS(port string) {
	log.Println("Listening for incoming ws on " + port)

	connectWs := http.NewServeMux()
	connectWs.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received ws connection from: " + r.RemoteAddr)

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Client connected: " + ws.RemoteAddr().String())

		wsClients[ws] = true
		go monitorWS(ws)
	})
	http.ListenAndServe(port, connectWs)
}

func monitorWS(conn *websocket.Conn) {
	for {
		if _, _, err := conn.NextReader(); err != nil {
			delete(wsClients, conn)
			conn.Close()
			log.Println("Client disconnected: " + conn.RemoteAddr().String())
			break
		}
	}
}

func consumeStream(stream <-chan []byte, done <-chan bool, recording *os.File) {
	for {
		select {
		case data := <-stream:
			if recording != nil {
				go writeToFile(data, recording)
			}

			for conn := range wsClients {
				if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					log.Println(err.Error())
				}
			}
		case <-done:
			for conn := range wsClients {
				conn.Close()
				delete(wsClients, conn)
			}
		}
	}
}

func writeToFile(data []byte, file *os.File) {
	if file != nil {
		file.Write(data)
	}
}
