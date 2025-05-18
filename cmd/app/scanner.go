package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const timeout = 500 * time.Millisecond

func scanPort(host string, port int, result chan <- int) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		result <- 0
		return
	}
}

func main() {
	if (len(os.Args)) < 2 {
		fmt.Println("Usage: gort <host>")
	}

	host := os.Args[1]
	startPort := 1
	endPort := 655535
	
	fmt.Printf("Scanning %v for open ports", host)
	result := make(chan int)
	var openPorts []int

	for port := startPort; port <= endPort; port++ {
		scanPort(host, port, result)
	}
}
