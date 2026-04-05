package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit


func main() {
	fmt.Println("=== Program starting... ===")

	// Setup listener
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {
        // Block until we receive an incoming connection
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
            //continue
        }

        // Handle client connection
        handleClient(conn)
		break
    }
}


func handleClient(conn net.Conn) {
	// Ensure we close the connection after we're done
    defer conn.Close()

	pongStr := []byte("+PONG\r\n")
	buf := make([]byte, 1024)

	for {
		bytesReceived, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}

		// Ignore empty commands
		if bytesReceived == 0 {
			continue
		}

		_, err = conn.Write(pongStr)
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			break
		}
	}
}