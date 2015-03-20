package main

type room struct {
	// forward: A queue for incoming messages that
	// should be forwarded to other clients
	forward chan []byte
}
