package main

import (
	"fmt"
	"net"
	"time"
	"os"
)

const timeout = 500 * time.Millisecond

func scanPort(host string, port int, results chan <- int) {
	
	address := fmt.Sprintf("%s:%d", host, port)
	connection, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		results <- 0
		return
	}
	
	connection.Close()
	results <- port
}

func main() {
	if (len(os.Args)) < 2 {
		fmt.Println("Usage: gort <host>")
	}

	host := os.Args[1]
	startPort := 1
	endPort := 65000

	fmt.Printf("Scanning %v for open ports...\n", host)

	results := make(chan int)
	var openPort []int

	for port := startPort; port <= endPort; port++ {
		go scanPort(host, port, results)
	}

	for i := startPort; i <= endPort; i++ {
		port := <-results
		if port != 0 {
			openPort = append(openPort, port)
		}
	}

	fmt.Printf("\nThe open ports in the host are: \n")
	for _, port := range openPort {
		fmt.Printf(" - Port %d is open\n", port)
	}
}
