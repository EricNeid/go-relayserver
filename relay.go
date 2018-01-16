package main

func relayStreamToWSClients(stream <-chan *[]byte, clients <-chan *wsClient) {
	go func() {
		newClient := <-clients
		go relayStreamToWS(stream, newClient)
	}()
}

func relayStreamToWS(stream <-chan *[]byte, client *wsClient) {

}
