package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type RESPValue interface {
	respValue()
	Serialize() []byte
}

type RESPArray struct {
	elements []RESPValue
}

type RESPBulkString struct {
	length int
	data   []byte
}

type RESPSimpleString struct {
	value string
}

// respValue() implementations
func (a RESPArray) respValue()        {}
func (b RESPBulkString) respValue()   {}
func (s RESPSimpleString) respValue() {}

// String() implementations
func (b RESPBulkString) String() string {
	return string(b.data)
}

func (s RESPSimpleString) String() string {
	return string(s.value)
}

// Serialize() implementations
func (a RESPArray) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(a.elements)))
	buf.WriteString("\r\n")
	for _, elem := range a.elements {
		buf.Write(elem.Serialize())
	}
	return buf.Bytes()
}

func (b RESPBulkString) Serialize() []byte {
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

func (s RESPSimpleString) Serialize() []byte {
	buf := make([]byte, 0, len(s.value)+3)
	buf = append(buf, '+')
	buf = append(buf, s.value...)
	buf = append(buf, '\r', '\n')
	return buf
}

// Parser Functions
func ParseResp(reader *bufio.Reader) (RESPValue, error) {
	firstByte, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	switch firstByte[0] {
	case '*':
		return parseRespArray(reader)
	case '$':
		return parseRespBulkString(reader)
	case '+':
		return parseRespSimpleString(reader)
	default:
		return nil, fmt.Errorf("unknown RESP type: %c", firstByte[0])
	}
}

func parseRespArray(reader *bufio.Reader) (RESPArray, error) {
	// Arrays:
	// payload[0] == "*"
	// *<number-of-elements>\r\n<element-1>...<element-n>
	// ["hello", "world"] -> *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	// *-1\r\n -> Null

	arrayHeader, err := reader.ReadBytes('\n')
	if err != nil {
		return RESPArray{}, err
	}

	if arrayHeader[0] != '*' {
		return RESPArray{}, fmt.Errorf("Invalid prefix for RESP array: %c", arrayHeader[0])
	}

	arrayCount, err := strconv.Atoi(string(arrayHeader[1 : len(arrayHeader)-2]))
	if err != nil {
		return RESPArray{}, err
	}

	if arrayCount == -1 {
		return RESPArray{}, nil
	}

	arrayElements := make([]RESPValue, 0, arrayCount)
	for range arrayCount {
		element, err := ParseResp(reader)
		if err != nil {
			return RESPArray{}, err
		}
		arrayElements = append(arrayElements, element)
	}

	return RESPArray{elements: arrayElements}, nil
}

func parseRespBulkString(reader *bufio.Reader) (RESPBulkString, error) {
	// Bulk Strings:
	// payload[0] == "$"
	// $<length>\r\n<data>\r\n
	// $11\r\nhello world\r\n
	// $0\r\n\r\n -> ""
	// $-1\r\n -> Null

	bulkHeader, err := reader.ReadBytes('\n')
	if err != nil {
		return RESPBulkString{}, err
	}

	if bulkHeader[0] != '$' {
		return RESPBulkString{}, fmt.Errorf("Invalid prefix for RESP bulk string: %c", bulkHeader[0])
	}

	bulkLength, err := strconv.Atoi(string(bulkHeader[1 : len(bulkHeader)-2]))
	if err != nil {
		return RESPBulkString{}, err
	}

	if bulkLength == -1 {
		return RESPBulkString{length: bulkLength}, nil
	}

	bulkData := make([]byte, bulkLength)
	_, err = io.ReadFull(reader, bulkData)
	if err != nil {
		return RESPBulkString{}, err
	}

	if _, err := reader.Discard(2); err != nil {
		return RESPBulkString{}, err
	}

	return RESPBulkString{length: bulkLength, data: bulkData}, nil
}

func parseRespSimpleString(reader *bufio.Reader) (RESPSimpleString, error) {
	// Simple Strings:
	// payload[0] == "+"
	// +<data>\r\n
	simpleStringPayload, err := reader.ReadBytes('\n')
	if err != nil {
		return RESPSimpleString{}, err
	}

	if simpleStringPayload[0] != '+' {
		return RESPSimpleString{}, fmt.Errorf("Invalid prefix for RESP simple string: %c", simpleStringPayload[0])
	}

	value := simpleStringPayload[1 : len(simpleStringPayload)-2]
	return RESPSimpleString{value: string(value)}, nil
}

func BulkStringsFromArray(respValue RESPValue) ([]RESPBulkString, bool) {
	array, ok := respValue.(RESPArray)
	if !ok {
		return nil, false
	}

	bulkStrings := make([]RESPBulkString, 0, len(array.elements))
	for _, element := range array.elements {
		element, ok := element.(RESPBulkString)
		if !ok {
			return nil, false
		}
		bulkStrings = append(bulkStrings, element)
	}
	return bulkStrings, true
}
