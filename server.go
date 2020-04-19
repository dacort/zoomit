package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

	code, hasCode := r.URL.Query()["code"]
	fmt.Println(code, hasCode)
	error, hasError := r.URL.Query()["error"]
	fmt.Println(error, hasError)
}

func makeHello(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

		code, _ := r.URL.Query()["code"]
		// fmt.Println(code, hasCode)
		// error, hasError := r.URL.Query()["error"]
		// fmt.Println(error, hasError)

		logger <- code[0]
	}
}

func serveOnARandomPort(responseChan chan string) (srv *http.Server, port string) {
	srv = &http.Server{}
	http.HandleFunc("/", makeHello(responseChan))
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Could not find a port for OAuth2 server to run on: %v", err)
	}
	port = fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
	fmt.Println("Using port:", port)
	go srv.Serve(listener)

	return srv, port
}

func getPortAndWait() (port string, waiter chan string) {
	// WE actually create a second channel here.
	// WE spin up the server, send back the HTTP port,
	// and also spin up a goroutine that waits for the server response.

	// When the repsonse comes back, we shutdown the server and then send
	// the code back on the channel.
	internalResponse := make(chan string)
	waiter = make(chan string)

	fmt.Println("Starting...")
	srv, port := serveOnARandomPort(internalResponse)

	go func(c chan string, s *http.Server) {
		fmt.Println("Waiting for OAuth response")
		msg := <-c
		fmt.Println("OAuth response received, shutting myself down...")
		s.Shutdown(context.Background())
		waiter <- msg
		fmt.Println("Grumble, grumble...")
	}(internalResponse, srv)

	return port, waiter
}

func serveAndWait() (string, error) {
	response := make(chan string)

	fmt.Println("Starting...")
	srv, _ := serveOnARandomPort(response)

	msg := <-response
	srv.Shutdown(context.Background())
	return msg, nil
}

// func main() {
// 	msg, _ := serveAndWait()
// 	fmt.Println(msg)
// 	// response := make(chan string)
// 	// serveOnARandomPort(response)

// 	// fmt.Println("Starting...")
// 	// msg := <-response
// 	// fmt.Println(msg)
// 	// // srv.Shutdown(context.Background())
// }
