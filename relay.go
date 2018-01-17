package main

import (
	"log"
)

// relayStreamToWSClients waits for clients to connect and relays the given stream to
// connected websocket clients. If a client disconnects, it does no longer receives the stream.
func relayStreamToWSClients(stream <-chan *[]byte, clients <-chan *wsClient) {
	connectedClients := make(map[*wsClient]bool)
	go func() {
		for {
			// wait for clients to connect
			newClient := <-clients
			log.Println("WS client connected: " + newClient.remoteAddress)
			// start goroutine to monitor the connection
			go func() {
				<-newClient.isClosed
				delete(connectedClients, newClient)
				log.Println("WS client disconnected: " + newClient.remoteAddress)
			}()
			connectedClients[newClient] = true
		}
	}()
	go func() {
		for {
			data := <-stream
			for client := range connectedClients {
				client.writeStream <- data
			}
		}
	}()
}
