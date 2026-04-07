package main

import (
	"fmt"
	"strings"
	"sync"
)

type Server struct {
	mu    sync.RWMutex
	store map[string][]byte
}

func NewServer() *Server {
	return &Server{
		store: make(map[string][]byte),
	}
}

type commandHandler func(s *Server, args []RESPBulkString) (RESPValue, error)

var commands = map[string]commandHandler{
	"PING": (*Server).handleCmdPing,
	"ECHO": (*Server).handleCmdEcho,
	"GET":  (*Server).handleCmdGet,
	"SET":  (*Server).handleCmdSet,
}

func (s *Server) HandleRequest(respRequest RESPValue) (RESPValue, error) {
	// "Clients send commands to a Redis server as an array of bulk strings.
	// The first (and sometimes also the second) bulk string in the array is the command's name.
	// Subsequent elements of the array are the arguments for the command."

	// "The server replies with a RESP type.
	// The reply's type is determined by the command's implementation
	// and possibly by the client's protocol version."

	bulkStrings, ok := BulkStringsFromArray(respRequest)
	if !ok {
		return nil, fmt.Errorf("Could not parse request as array of bulk strings")
	}

	command := strings.ToUpper(bulkStrings[0].String())
	args := bulkStrings[1:]
	handler, ok := commands[command]
	if !ok {
		return nil, fmt.Errorf("unknown command: %s", command)
	}
	return handler(s, args)
}

func (s *Server) handleCmdSet(args []RESPBulkString) (RESPValue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[args[0].String()] = args[1].data
	return NewRESPSimpleString("OK"), nil
}

func (s *Server) handleCmdGet(args []RESPBulkString) (RESPValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.store[args[0].String()]
	if !ok {
		return NewNullRESPBulkString(), nil
	}
	return NewRESPBulkString(data), nil
}

func (s *Server) handleCmdEcho(args []RESPBulkString) (RESPValue, error) {
	// ECHO only handles a single arg
	return args[0], nil
}

func (s *Server) handleCmdPing(_ []RESPBulkString) (RESPValue, error) {
	return NewRESPSimpleString("PONG"), nil
}
