package main

import (
	"log"
)

func relayStreamToWSClients(stream <-chan *[]byte, clients <-chan *wsClient) {
	connectedClients := make(map[*wsClient]bool)
	go func() {
		for {
			// wait for clients to connect
			newClient := <-clients
			log.Println("Client connected: " + newClient.remoteAddress)
			// start goroutine to monitor the connection
			go func() {
				<-newClient.isClosed
				delete(connectedClients, newClient)
				log.Println("Client disconnected: " + newClient.remoteAddress)
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
