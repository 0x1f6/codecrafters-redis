package redis

import "github.com/0x1f6/codecrafters-redis/internal/resp"

func (d *Dispatcher) handleCmdRpush(args []resp.BulkString) (resp.RESPValue, error) {
	key, values, err := parseRpushArgs(args)
	if err != nil {
		return nil, err
	}
	count, ok := d.store.Rpush(key, values...)
	if !ok {
		return nil, errorf("RPUSH failed")
	}
	return resp.NewInteger(count), nil
}

func parseRpushArgs(args []resp.BulkString) (string, [][]byte, error) {
	if len(args) < 2 {
		return "", nil, errorf("RPUSH requires at least 2 arguments")
	}
	key := string(args[0].Data())
	values := make([][]byte, 0, len(args)-1)
	for _, value := range args[1:] {
		values = append(values, value.Data())
	}

	return key, values, nil
}