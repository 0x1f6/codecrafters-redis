package redis

import (
	"fmt"
	"strings"

	"github.com/0x1f6/codecrafters-redis/internal/resp"
)

type Dispatcher struct {
	store *Store
}

type commandHandler func(d *Dispatcher, args []resp.BulkString) (resp.RESPValue, error)

func New() *Dispatcher {
	return &Dispatcher{
		store: NewStore(),
	}
}

var commands = map[string]commandHandler{
	"PING": (*Dispatcher).handleCmdPing,
	"ECHO": (*Dispatcher).handleCmdEcho,
	"GET":  (*Dispatcher).handleCmdGet,
	"SET":  (*Dispatcher).handleCmdSet,
}

func errorf(format string, args ...any) resp.SimpleError {
	return resp.NewSimpleError("ERR " + fmt.Sprintf(format, args...))
}

func (d *Dispatcher) HandleRequest(respRequest resp.RESPValue) (resp.RESPValue, error) {
	// "Clients send commands to a Redis server as an array of bulk strings.
	// The first (and sometimes also the second) bulk string in the array is the command's name.
	// Subsequent elements of the array are the arguments for the command."

	// "The server replies with a RESP type.
	// The reply's type is determined by the command's implementation
	// and possibly by the client's protocol version."

	requestArray, ok := respRequest.(resp.Array)
	if !ok {
		return nil, errorf("expected request to be an array of bulk strings")
	}

	bulkStrings, ok := requestArray.AsBulkStrings()
	if !ok {
		return nil, errorf("could not parse request as array of bulk strings")
	}

	command := strings.ToUpper(bulkStrings[0].String())
	args := bulkStrings[1:]
	handler, ok := commands[command]
	if !ok {
		return nil, errorf("unknown command '%s'", command)
	}
	return handler(d, args)
}
