package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type BulkString struct {
	length int
	data   []byte
}

func (b BulkString) respValue() {}

func (b BulkString) String() string {
	return string(b.data)
}

func (b BulkString) Serialize() []byte {
	if b.length == -1 {
		return []byte("$-1\r\n")
	}
	buf := make([]byte, 0, len(b.data)+16)
	buf = append(buf, '$')
	buf = strconv.AppendInt(buf, int64(b.length), 10)
	buf = append(buf, '\r', '\n')
	buf = append(buf, b.data...)
	buf = append(buf, '\r', '\n')
	return buf
}

func NewBulkString(data []byte) BulkString {
	return BulkString{length: len(data), data: data}
}

func NewNullBulkString() BulkString {
	return BulkString{length: -1, data: nil}
}

func (b BulkString) IsNull() bool {
	return b.length == -1
}

func (b BulkString) Data() []byte {
	return b.data
}

func (b BulkString) Length() int {
	return b.length
}

func parseBulkString(reader *bufio.Reader) (BulkString, error) {
	// Bulk Strings:
	// payload[0] == "$"
	// $<length>\r\n<data>\r\n
	// $11\r\nhello world\r\n
	// $0\r\n\r\n -> ""
	// $-1\r\n -> Null

	bulkHeader, err := reader.ReadBytes('\n')
	if err != nil {
		return BulkString{}, err
	}

	if bulkHeader[0] != '$' {
		return BulkString{}, fmt.Errorf("Invalid prefix for RESP bulk string: %c", bulkHeader[0])
	}

	bulkLength, err := strconv.Atoi(string(bulkHeader[1 : len(bulkHeader)-2]))
	if err != nil {
		return BulkString{}, err
	}

	if bulkLength == -1 {
		return BulkString{length: bulkLength}, nil
	}

	bulkData := make([]byte, bulkLength)
	_, err = io.ReadFull(reader, bulkData)
	if err != nil {
		return BulkString{}, err
	}

	if _, err := reader.Discard(2); err != nil {
		return BulkString{}, err
	}

	return BulkString{length: bulkLength, data: bulkData}, nil
}
