package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Server create a Proxy server on listenPort to forward traffic to target and targetPort
func Server(listenPort int, target string, targetPort int) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handle(target, targetPort, conn)
	}
}

func handle(target string, targetPort int, src net.Conn) {
	dst, err := net.Dial("tcp", fmt.Sprintf("%s:%d", target, targetPort))
	if err != nil {
		log.Fatalln("Unable to connect to target host")
	}
	defer dst.Close()

	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}
