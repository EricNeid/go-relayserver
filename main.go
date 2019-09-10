package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

const defaultPortStream = ":8000"
const defaultPortWS = ":9000"
const defaultSecret = "secret1234"
const defaultRecordToFile = false

const bufferSize = 8 * 1000 * 1024 // 8MB

var recordName = fmt.Sprintf("%s_record.mpeg", time.Now().Format("20060102_1504"))

type config struct {
	portStream   string
	portWS       string
	secretStream string
	record       bool
	printHelp    bool
}

func main() {
	// read arguments from cli
	config := readCmdArguments()
	if config.printHelp {
		log.Println("Usage: ")
		log.Println("go-relayserver.exe optional: -port-stream <port> -port-ws <port> -s <secret>")
		return
	}

	// start listening
	done := make(chan bool, 1)
	_, stream := waitForStream(config.portStream, config.secretStream, done)
	if config.record {
		log.Println("Recording stream to " + recordName)
		log.Println("Warning: Recording stream may decrease performance and should be used for testing only")
		stream = recordStream(stream, "recorded", recordName)
	}
	clients := waitForWSClients(config.portWS)
	relayStreamToWSClients(stream, clients)

	// wait for interrupt to shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	log.Println("Shuting down...")
	done <- true
}

func readCmdArguments() config {
	help := flag.Bool("h", false, "print help")

	portStream := flag.String("port-stream", defaultPortStream, "Port to listen for stream")
	portWS := flag.String("port-ws", defaultPortWS, "Port to listen for websockets")
	secretStream := flag.String("s", defaultSecret, "Secure stream with this password")
	record := flag.Bool("record", defaultRecordToFile, "Record received stream to local file")

	flag.Parse()

	return config{
		portStream:   normalizePort(*portStream),
		portWS:       normalizePort(*portWS),
		record:       *record,
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
