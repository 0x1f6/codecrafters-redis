package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func ServeForever() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379: ", err.Error())
		return
	}

	defer listener.Close()

	s := NewServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {

		request, err := ParseResp(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error parsing request:", err)
			break
		}

		response, err := s.HandleRequest(request)
		if err != nil {
			fmt.Println("Error handling request:", err)
			break
		}

		_, err = conn.Write(response.Serialize())
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			break
		}
	}
}
