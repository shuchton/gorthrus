package echo

import (
	"fmt"
	"io"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()

	// Copy data via io.Copy
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}

// Server establishes a TCP Echo server on the given port
func Server(port string) {
	// Bind to TCP port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln("Unable to bind to port", port)
	}
	log.Printf("Listening on 0.0.0.0:%s\n", port)
	
	for {
		// Wait for connection
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		// Handle the connection with concurrency
		go echo(conn)
	}
}
