package main

import (
	"flag"
	"fmt"
	"strings"
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

	stream := waitForStream(config.portStream, config.secretStream)
	if recordToFile {
		fmt.Println("Recording stream to " + recordName)
		fmt.Println("Warning: Recording stream may decrease performance and shoudl be used for testing only")
		stream = recordStream(stream, recordName)
	}

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
		portStream:   normalizePort(*portStream),
		portWS:       normalizePort(*portWS),
		secretStream: *secretStream,
		printHelp:    *help,
	}
}

func normalizePort(port string) string {
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}
	return port
}
