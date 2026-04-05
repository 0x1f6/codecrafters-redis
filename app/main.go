package main

import (
	"fmt"
	"io"
	"net"
	"os"
)


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
            continue
        }

        // Handle client connection
        go handleClient(conn)
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
			if err != io.EOF {
				fmt.Println("Error reading from connection: ", err.Error())
			}
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