package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)


var CommonPorts = map[string]string{
	"1/tcp":     "tcpmux",      // TCP Port Service Multiplexer
	"7/tcp":     "echo",        // Echo Protocol
	"7/udp":     "echo",
	"9/tcp":     "discard",     // Discard Protocol
	"9/udp":     "discard",
	"13/tcp":    "daytime",     // Daytime Protocol
	"13/udp":    "daytime",
	"17/tcp":    "qotd",        // Quote of the Day
	"17/udp":    "qotd",
	"19/tcp":    "chargen",     // Character Generator Protocol
	"19/udp":    "chargen",
	"20/tcp":    "ftp-data",    // FTP Data Transfer
	"21/tcp":    "ftp",         // FTP Control
	"22/tcp":    "ssh",         // Secure Shell
	"23/tcp":    "telnet",      // Telnet
	"25/tcp":    "smtp",        // Simple Mail Transfer Protocol
	"37/tcp":    "time",        // Time Protocol
	"37/udp":    "time",
	"43/tcp":    "whois",       // WHOIS Protocol
	"53/tcp":    "dns",         // Domain Name System
	"53/udp":    "dns",
	"67/udp":    "dhcp",        // DHCP Server
	"68/udp":    "dhcp",        // DHCP Client
	"69/udp":    "tftp",        // Trivial File Transfer Protocol
	"70/tcp":    "gopher",      // Gopher Protocol
	"79/tcp":    "finger",      // Finger Protocol
	"80/tcp":    "http",        // Hypertext Transfer Protocol
	"88/tcp":    "kerberos",    // Kerberos Authentication
	"88/udp":    "kerberos",
	"110/tcp":   "pop3",        // Post Office Protocol v3
	"111/tcp":   "rpcbind",     // RPC Portmapper
	"111/udp":   "rpcbind",
	"119/tcp":   "nntp",        // Network News Transfer Protocol
	"123/udp":   "ntp",         // Network Time Protocol
	"135/tcp":   "msrpc",       // Microsoft RPC
	"137/udp":   "netbios-ns",  // NetBIOS Name Service
	"138/udp":   "netbios-dgm", // NetBIOS Datagram Service
	"139/tcp":   "netbios-ssn", // NetBIOS Session Service
	"143/tcp":   "imap",        // Internet Message Access Protocol
	"161/udp":   "snmp",        // Simple Network Management Protocol
	"162/udp":   "snmptrap",    // SNMP Trap
	"179/tcp":   "bgp",         // Border Gateway Protocol
	"194/tcp":   "irc",         // Internet Relay Chat
	"389/tcp":   "ldap",        // Lightweight Directory Access Protocol
	"443/tcp":   "https",       // HTTP Secure
	"445/tcp":   "smb",         // Server Message Block (Microsoft SMB)
	"465/tcp":   "smtps",       // SMTP Secure
	"514/udp":   "syslog",      // Syslog
	"515/tcp":   "lpd",         // Line Printer Daemon
	"520/udp":   "rip",         // Routing Information Protocol
	"587/tcp":   "submission",  // SMTP Submission
	"636/tcp":   "ldaps",       // LDAP Secure
	"993/tcp":   "imaps",       // IMAP Secure
	"995/tcp":   "pop3s",       // POP3 Secure
	"1080/tcp":  "socks",       // SOCKS Proxy
	"1433/tcp":  "mssql",       // Microsoft SQL Server
	"1521/tcp":  "oracle",      // Oracle Database
	"1723/tcp":  "pptp",        // Point-to-Point Tunneling Protocol
	"2049/tcp":  "nfs",         // Network File System
	"2049/udp":  "nfs",
	"3306/tcp":  "mysql",       // MySQL Database
	"3389/tcp":  "rdp",         // Remote Desktop Protocol
	"5432/tcp":  "postgresql",  // PostgreSQL Database
	"5900/tcp":  "vnc",         // Virtual Network Computing
	"6379/tcp":  "redis",       // Redis Database
	"8080/tcp":  "http-alt",    // Alternative HTTP (often used for web servers)
	"8443/tcp":  "https-alt",   // Alternative HTTPS
	"27017/tcp": "mongodb",     // MongoDB Database
}

func measureLatency(host string) (time.Duration, error) {

	start := time.Now()
	connection, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", host), 2 * time.Second)
	if err == nil {
		connection.Close()
		return time.Since(start), nil
	}

	start = time.Now()
	connection, err = net.DialTimeout("tcp", fmt.Sprintf("%s:443", host), 2 * time.Second)
	if err == nil {
		connection.Close()
		return time.Since(start), nil
	}

	return 500 * time.Millisecond, fmt.Errorf("could not measure latency")
}

func calculateTimeout(rtt time.Duration) time.Duration {
	const (
		minTimeout = 100 * time.Millisecond
		maxTimeout = 1000 * time.Millisecond
		multiplier = 3
	)

	timeout := rtt * multiplier
	if timeout < minTimeout {
		return minTimeout
	}
	if timeout > maxTimeout {
		return maxTimeout
	}

	return timeout
}

func getService(port int, protocol string) string {
	key := fmt.Sprintf("%d/%s", port, protocol)
	if services, exists := CommonPorts[key]; exists {
		return services
	}
	return "unknown"
}

func scanPort(host string, port int, results chan <- int, timeout time.Duration) {
	
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
	var openPort []int

	// Dynamic timeout allocation
	rtt, err := measureLatency(host)
	timeout := calculateTimeout(rtt)

	fmt.Printf("%v %v", rtt, timeout)
	for port := startPort; port <= endPort; port++ {
		go scanPort(host, port, results, timeout)
	}

	for i := startPort; i <= endPort; i++ {
		port := <-results
		if port != 0 {
			openPort = append(openPort, port)
		}
	}

	fmt.Printf("\nThe open ports in the host are: \n")
	fmt.Printf("Port\tStatus\tService\n")
	for _, port := range openPort {
		protocol := "tcp"
		service := getService(port, protocol)
		fmt.Printf("%v/%v\topen\t%v\n", port, protocol, service)
	}
}
