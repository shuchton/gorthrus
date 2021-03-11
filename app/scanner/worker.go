package scan

import (
	"fmt"
	"net"
)

func worker(network, site string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", site, p)
		conn, err := net.Dial(network, address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}
