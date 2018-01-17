package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const recordToFile = false

const defaultPortStream = ":8081"
const defaultPortWS = ":8082"
const bufferSize = 8 * 1000 * 1024 // 8MB

var recordName = fmt.Sprintf("%s_record.mpeg", time.Now().Format("20060102_1504"))

type config struct {
	portStream   string
	portWS       string
	secretStream string
	printHelp    bool
}

func main() {
	config := readCmdArguments()
	if config.printHelp {
		fmt.Println("Usage: ")
		fmt.Println("go-relayserver.exe optional: -port-stream <port> -port-ws <port> -s <secret>")
		return
	}

	done := make(chan bool)

	// var recording *os.File
	// if recordToFile {
	// 	f, err := os.Create(recordName)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	} else {
	// 		defer f.Close()
	// 		recording = f
	// 	}
	// }

	stream := waitForStream(config.portStream, config.secretStream)
	clients := waitForWSClients(config.portWS)
	relayStreamToWSClients(stream, clients)

	fmt.Println("Relay started, hit Enter-key to close")

	fmt.Scanln()
	done <- true

	fmt.Println("Shuting down...")
}

func readCmdArguments() config {
	help := flag.Bool("h", false, "print help")

	portStream := flag.String("port-stream", defaultPortStream, "Port to listen for stream")
	portWS := flag.String("port-ws", defaultPortWS, "Port to listen for websockets")
	secretStream := flag.String("s", "", "Secure stream with this password")

	flag.Parse()

	return config{
		portStream:   *portStream,
		portWS:       *portWS,
		secretStream: *secretStream,
		printHelp:    *help,
	}
}

func waitForStream(port string, secret string) <-chan *[]byte {
	log.Println("Listening for incoming stream on " + port)

	stream := make(chan *[]byte)
	go func() {
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
				stream <- &chunk
			}
			log.Println("Stream closed")
		})
		http.ListenAndServe(port, streamReader)
	}()
	return stream
}
