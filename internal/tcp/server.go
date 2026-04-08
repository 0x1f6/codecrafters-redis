package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/0x1f6/codecrafters-redis/internal/resp"
)

type Server struct {
	addr    string
	handler RequestHandler
}

type RequestHandler interface {
	HandleRequest(request resp.RESPValue) (resp.RESPValue, error)
}

func NewServer(addr string, handler RequestHandler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

func (s *Server) ServeForever() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection: ", err.Error())
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		request, err := resp.Parse(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error parsing request:", err)
			break
		}

		response, err := s.handler.HandleRequest(request)
		if err != nil {
			simpleErr, ok := err.(resp.SimpleError)
			if !ok {
				fmt.Println("Error handling request:", err)
				break
			}
			response = simpleErr
		}

		_, err = conn.Write(response.Serialize())
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			break
		}
	}
}
