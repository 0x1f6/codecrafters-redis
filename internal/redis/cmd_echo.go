package redis

import "github.com/0x1f6/codecrafters-redis/internal/resp"

func (d *Dispatcher) handleCmdEcho(args []resp.BulkString) (resp.RESPValue, error) {
	// ECHO only handles a single arg
	return args[0], nil
}
