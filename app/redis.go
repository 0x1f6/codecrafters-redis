package main

import (
	"fmt"
	"strings"
)

func HandleRequest(respRequest RESPValue) ([]byte, error) {
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

	command, args := strings.ToLower(bulkStrings[0].String()), bulkStrings[1:]

	switch command {
	case "ping":
		return handleCmdPing(args)
	case "echo":
		return handleCmdEcho(args)
	default:
		return nil, fmt.Errorf("unknown command: %s", command)

	}
}

func handleCmdEcho(args []RESPBulkString) ([]byte, error) {
	// ECHO only handles a single arg
	return args[0].Serialize(), nil
}

func handleCmdPing(_ []RESPBulkString) ([]byte, error) {
	return []byte("+PONG\r\n"), nil
}
