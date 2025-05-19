package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"strings"
)

const timeout = 500 * time.Millisecond

func getService(port int, protocol string) (string, error) {

	file, err := os.Open("/etc/services")
	if err != nil {
		return "unknown", err // In case of windows
	}

	scanner := bufio.NewScanner(file)
	portProtocol := fmt.Sprintf("%s/%d", protocol, port)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		args := strings.Fields(line)
		if (len(args)) < 2 {
			continue
		}
		if args[1] == portProtocol {
			return args[0], nil
		}

	}
	return "unknown", scanner.Err()
}

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
	if (len(os.Args)) < 3 {
		fmt.Println("Usage: gort <host>")
	}

	host := os.Args[1]
	startPort := 0
	endPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Scanning %v for open ports...\n", host)

	results := make(chan int)
	var openPort []string

	for port := startPort; port <= endPort; port++ {
		go scanPort(host, port, results)
	}

	for i := startPort; i <= endPort; i++ {
		port := <-results
		if port != 0 {
			openPort = append(openPort, fmt.Sprintf("%d/tcp", port))
		}
	}

	fmt.Printf("\nThe open ports in the host are: \n")
	fmt.Printf("Port\tStatus\tService\n")
	for _, port := range openPort {
		fmt.Printf("%v\topen\t\n", port)
	}
}
