package redis

import "github.com/0x1f6/codecrafters-redis/internal/resp"

func (d *Dispatcher) handleCmdGet(args []resp.BulkString) (resp.RESPValue, error) {
	data, ok := d.store.Get(args[0].String())
	if !ok {
		return resp.NewNullBulkString(), nil
	}
	return resp.NewBulkString(data), nil
}
