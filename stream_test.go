package main




/*
func TestWaitForStream(t *testing.T) {
	// arrange
	done := make(chan bool, 1)
	server, stream := waitForStream(":8990", "test", done)

	// action
	go func() {
		err := sendData(":8990", "test", "Hallo, Welt")
		ok(t, err)
	}()

	// verify
	received := <-stream
	assert(t, received != nil, "Received data is null")
	equals(t, "Hallo, Welt", string(*received))

	// action
	go func() {
		err := sendData(":8990", "test", "Hallo, Welt 2")
		ok(t, err)
	}()

	// verify
	received = <-stream
	assert(t, received != nil, "Received data is null")
	equals(t, "Hallo, Welt 2", string(*received))

	// cleanup
	done <- true
	server.Shutdown(context.Background())
}

func TestRecordStream(t *testing.T) {
	// arrange
	done := make(chan bool, 1)
	os.Remove("testdata/recorded-sample.txt")
	server, stream := waitForStream(":8990", "test", done)
	go func() {
		streamRecorded := recordStream(stream, "testdata", "recorded-sample.txt")
		for {
			<-streamRecorded
		}
	}()

	// action
	err := sendData(":8990", "test", "Hallo, Welt")

	// verify
	ok(t, err)
	recorded, err := ioutil.ReadFile("testdata/recorded-sample.txt")
	ok(t, err)
	equals(t, "Hallo, Welt", string(recorded))

	// cleanup
	server.Shutdown(context.Background())
	done <- true
}

*/