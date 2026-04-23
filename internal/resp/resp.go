package resp

import (
	"bufio"
	"fmt"
)

type RESPValue interface {
	respValue()
	Serialize() []byte
}

func Parse(reader *bufio.Reader) (RESPValue, error) {
	firstByte, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	switch firstByte[0] {
	case '*':
		return parseArray(reader)
	case '$':
		return parseBulkString(reader)
	case '+':
		return parseSimpleString(reader)
	case '-':
		return parseSimpleError(reader)
	case ':':
		return parseInteger(reader)
	default:
		return nil, fmt.Errorf("unknown RESP type: %c", firstByte[0])
	}
}
