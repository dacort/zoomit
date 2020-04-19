package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

func handleOauthResponse(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You may now close this window! 🎉")

		// Parse the OAuth response code
		code, _ := r.URL.Query()["code"]

		// Send it back to the waiting channel
		logger <- code[0]
	}
}

func serveOnARandomPort(responseChan chan string) (srv *http.Server, port string) {
	srv = &http.Server{}
	http.HandleFunc("/", handleOauthResponse(responseChan))
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Could not find a port for OAuth2 server to run on: %v", err)
	}
	port = fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
	go srv.Serve(listener)

	return srv, port
}

func getPortAndWait() (port string, waiter chan string) {
	// We actually create a second channel here.
	// We spin up the server, send back the HTTP port,
	// and also spin up a goroutine that waits for the server response.

	// When the response comes back, we shutdown the server and then send
	// the code back on the channel.
	internalResponse := make(chan string)
	waiter = make(chan string)

	log.Println("Starting mini HTTP server...")
	srv, port := serveOnARandomPort(internalResponse)

	go func(c chan string, s *http.Server) {
		log.Printf("Waiting for OAuth response on :%s\n", port)
		msg := <-c
		s.Shutdown(context.Background())
		waiter <- msg
	}(internalResponse, srv)

	return port, waiter
}
