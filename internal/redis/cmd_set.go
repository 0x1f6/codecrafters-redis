package redis

import (
	"strconv"
	"strings"
	"time"

	"github.com/0x1f6/codecrafters-redis/internal/resp"
)

type setOptions struct {
	expiresAt time.Time
}

func (d *Dispatcher) handleCmdSet(args []resp.BulkString) (resp.RESPValue, error) {
	key, value, opts, err := parseSetArgs(args)
	if err != nil {
		return nil, err
	}
	d.store.Set(key, value, opts)
	return resp.NewSimpleString("OK"), nil
}

func parseSetArgs(args []resp.BulkString) (string, []byte, setOptions, error) {
	if len(args) < 2 {
		return "", nil, setOptions{}, errorf("SET requires at least 2 arguments")
	}
	key := string(args[0].Data())
	value := args[1].Data()

	opts, err := parseSetOptions(args[2:])
	if err != nil {
		return "", nil, setOptions{}, err
	}
	return key, value, opts, nil
}

func parseSetOptions(args []resp.BulkString) (setOptions, error) {
	var expiresAt time.Time
	var err error

	for len(args) > 0 {
		switch strings.ToUpper(args[0].String()) {
		case "EX":
			if !expiresAt.IsZero() {
				return setOptions{}, errorf("multiple expire times in set")
			}
			expiresAt, args, err = parseExpiration(args, time.Second)
			if err != nil {
				return setOptions{}, err
			}
		case "PX":
			if !expiresAt.IsZero() {
				return setOptions{}, errorf("multiple expire times in set")
			}
			expiresAt, args, err = parseExpiration(args, time.Millisecond)
			if err != nil {
				return setOptions{}, err
			}
		default:
			return setOptions{}, errorf("invalid option in set: %s", args[0].String())
		}
	}
	return setOptions{expiresAt: expiresAt}, nil
}

func parseExpiration(args []resp.BulkString, unit time.Duration) (time.Time, []resp.BulkString, error) {
	if len(args) < 2 {
		return time.Time{}, nil, errorf("missing expire time in set")
	}

	val, err := strconv.Atoi(args[1].String())
	if err != nil {
		return time.Time{}, nil, errorf("invalid expire time in set")
	}
	if val < 1 {
		return time.Time{}, nil, errorf("invalid expire time in set")
	}
	expiresIn := time.Duration(val) * unit
	return time.Now().Add(expiresIn), args[2:], nil
}
