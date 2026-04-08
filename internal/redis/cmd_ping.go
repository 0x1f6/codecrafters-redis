package redis

import "github.com/0x1f6/codecrafters-redis/internal/resp"

func (d *Dispatcher) handleCmdPing(_ []resp.BulkString) (resp.RESPValue, error) {
	return resp.NewSimpleString("PONG"), nil
}
