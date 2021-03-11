package scan

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

// Scanner scans the network.
func Scanner(network, site, scanPorts string) ([]int, error) {
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	parsedPorts, err := parsePorts(scanPorts)
	if err != nil {
		return nil, err
	}

	for i := 0; i < cap(ports); i++ {
		go worker(network, site, ports, results)
	}

	go func() {
		for _, i := range parsedPorts {
			ports <- i
		}
	}()

	for i := 0; i < len(parsedPorts); i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	return openports, nil
}

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

// parsePorts takes a string input for a port, a port range [1-1024]
// or a list of ports 80,22 or a compination like 80,20-35,443 and generates
// an array of ports to scan.
func parsePorts(portRange string) ([]int, error) {
	addedPorts := make(map[int]struct{})

	const (
		min = 1
		max = uint32(1<<32 - 1)
	)

	// Remove all white space from the string
	_portRange := strings.ReplaceAll(portRange, " ", "")

	// Split the port range by commas
	portArray := strings.Split(_portRange, ",")

	for _, ports := range portArray {
		// Check is ports is a range
		portsArray := strings.Split(ports, "-")
		if len(portsArray) > 1 {
			start, err := strconv.Atoi(portsArray[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(portsArray[1])
			if err != nil {
				return nil, err
			}
			if start < min {
				start = min
			}
			if uint32(end) > max {
				end = int(max)
			}
			for i := start; i <= end; i++ {
				addedPorts[i] = struct{}{}
			}
		} else {
			port, err := strconv.Atoi(portsArray[0])
			if err != nil {
				return nil, err
			}
			addedPorts[port] = struct{}{}
		}
	}

	ret := make([]int, len(addedPorts))
	idx := 0
	for k := range addedPorts {
		ret[idx] = k
		idx++
	}

	sort.Ints(ret)

	return ret, nil
}
