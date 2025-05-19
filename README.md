# Gort - A Simple Port Scanner in Go

Gort is a lightweight, concurrent port scanner written in Go. It scans a specified host for open TCP ports and identifies the services running on those ports using a predefined map of common port-to-service mappings. The scanner is designed for educational purposes and network administration tasks.

## Features
- Concurrent TCP port scanning for fast results
- Service identification for common ports (e.g., 80/tcp -> http)
- Simple command-line interface
- Cross-platform (Linux, macOS, Windows)

## Installation

1. **Install Go**: Ensure you have Go (version 1.16 or later) installed. Download it from [golang.org](https://golang.org/dl/).

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/BitsExploited/gortscanner.git
   cd gortscanner
   go build cmd/app/scanner -o gort
   ```
3. **Run the scanner**
    ```bash
    gort <host> <range_of_port>
    ```
## Development Status

- [x] TCP port scanning
- [x] Service identification using common ports map
- [ ] Service indentification using client-side script
- [ ] UDP port scanning
- [ ] Custom start port option
- [ ] Output to file

