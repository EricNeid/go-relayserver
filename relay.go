package main

// relayStreamToWSClients waits for clients to connect and relays the given stream to
// connected websocket clients. If a client disconnects, it does no longer receives the stream.
func relayStreamToWSClients(stream <-chan *[]byte, clients <-chan *wsClient) {
	connectedClients := make(map[*wsClient]bool)
	go func() {
		for {
			// wait for clients to connect
			newClient := <-clients

			// start goroutine to monitor the connection for disconnect
			go func() {
				<-newClient.isClosed
				delete(connectedClients, newClient)
			}()
			connectedClients[newClient] = true
		}
	}()
	// read from stream and forward to each client
	go func() {
		for {
			data := <-stream
			for client := range connectedClients {
				client.writeStream <- data
			}
		}
	}()
}
