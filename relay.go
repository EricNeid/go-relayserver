package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const recordToFile = false

const defaultPortStream = ":8081"
const defaultPortWS = ":8082"
const bufferSize = 16 * 1000 * 1024 // 16MB

var wsClients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var recordName = fmt.Sprintf("%s_record.mpeg", time.Now().Format("20060102_1504"))

type input struct {
	portStream   string
	portWS       string
	secretStream string
	printHelp    bool
}

func main() {
	input := readCmdArguments()
	if input.printHelp {
		fmt.Println("Usage: ")
		fmt.Println("go-relayserver.exe optional: -port-stream <port> -port-ws <port> -s <secret>")
		return
	}

	stream := make(chan []byte)
	done := make(chan bool)

	var recording *os.File
	if recordToFile {
		f, err := os.Create(recordName)
		if err != nil {
			log.Println(err.Error())
		} else {
			defer f.Close()
			recording = f
		}
	}

	go consumeStream(stream, done, recording)
	go waitForStream(input.portStream, input.secretStream, stream)
	go waitForWS(input.portWS)

	fmt.Println("Relay started, hit Enter-key to close")

	fmt.Scanln()
	done <- true

	fmt.Println("Shuting down...")
}

func readCmdArguments() input {
	help := flag.Bool("h", false, "print help")

	portStream := flag.String("port-stream", defaultPortStream, "Port to listen for stream")
	portWS := flag.String("port-ws", defaultPortWS, "Port to listen for websockets")
	secretStream := flag.String("s", "", "Secure stream with this password")

	flag.Parse()

	return input{
		portStream:   *portStream,
		portWS:       *portWS,
		secretStream: *secretStream,
		printHelp:    *help,
	}
}

func waitForStream(port string, secret string, stream chan<- []byte) {
	log.Println("Listening for incoming stream on " + port)

	streamReader := http.NewServeMux()
	streamReader.HandleFunc("/"+secret, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received stream connection from: " + r.RemoteAddr)

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

// monitorWS waits for the websocket to disconnect.
// It closes the connection and removes the ws from the connected pool.
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

// consumeStream reads from stream channel and writes the []byte to each connected websocket.
// If done is set, it closes all connections.
// If given file for recording is not null, the stream is also written to the file.
func consumeStream(stream <-chan []byte, done <-chan bool, recording *os.File) {
	for {
		select {
		case data := <-stream:
			if recording != nil {
				// decouple file writing in
				writeToFile := func(data []byte, file *os.File) {
					if file != nil {
						file.Write(data)
					}
				}
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
