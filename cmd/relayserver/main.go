package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	relay "github.com/EricNeid/go-relayserver"
)

const defaultPortStream = ":8000"
const defaultPortWS = ":9000"
const defaultSecret = "secret1234"
const defaultRecordToFile = false

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
		log.Println("relayserver.exe optional: -port-stream <port> -port-ws <port> -s <secret>")
		return
	}

	// start endpoints for stream and websocket connections
	streamServer := relay.NewStreamServer(config.portStream, config.secretStream)
	streamServer.Routes()

	wsServer := relay.NewWebSocketServer(config.portWS)
	wsServer.Routes()

	go func() {
		streamServer.ListenAndServe()
	}()
	go func() {
		wsServer.ListenAndServe()
	}()

	var stream <-chan *[]byte
	stream = streamServer.InputStream
	clients := wsServer.IncomingClients
	if config.record {
		log.Println("Recording stream to " + recordName)
		log.Println("Warning: Recording stream may decrease performance and should be used for testing only")
		stream = relay.RecordStream(stream, "recorded", recordName)
	}
	relay.StreamToWSClients(stream, clients)

	// wait for interrupt to shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	log.Println("Shuting down...")
	// TODO cleanup
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
