package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

type Array struct {
	count    int
	elements []RESPValue
}

func (a Array) respValue() {}

func NewArray(elements []RESPValue) Array {
	return Array{count: len(elements), elements: elements}
}

func NewNullArray() Array {
	return Array{count: -1, elements: nil}
}

func (a Array) IsNull() bool {
	return a.count == -1
}

func (a Array) Elements() []RESPValue {
	return a.elements
}

func (a Array) Element(i int) RESPValue {
	return a.elements[i]
}

func (a Array) Count() int {
	return a.count
}

func (a Array) Serialize() []byte {
	if a.count == -1 {
		return []byte("*-1\r\n")
	}
	var buf bytes.Buffer
	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(a.count))
	buf.WriteString("\r\n")
	for _, elem := range a.elements {
		buf.Write(elem.Serialize())
	}
	return buf.Bytes()
}

func parseArray(r *bufio.Reader) (Array, error) {
	// Arrays:
	// payload[0] == "*"
	// *<number-of-elements>\r\n<element-1>...<element-n>
	// ["hello", "world"] -> *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	// *-1\r\n -> Null

	header, err := r.ReadBytes('\n')
	if err != nil {
		return Array{}, err
	}

	if header[0] != '*' {
		return Array{}, fmt.Errorf("Invalid prefix for RESP array: %c", header[0])
	}

	count, err := strconv.Atoi(string(header[1 : len(header)-2]))
	if err != nil {
		return Array{}, err
	}

	if count == -1 {
		return Array{}, nil
	}

	elements := make([]RESPValue, 0, count)
	for range count {
		element, err := Parse(r)
		if err != nil {
			return Array{}, err
		}
		elements = append(elements, element)
	}

	return Array{elements: elements}, nil
}

func (a Array) AsBulkStrings() ([]BulkString, bool) {
	bulkStrings := make([]BulkString, 0, len(a.elements))
	for _, element := range a.elements {
		element, ok := element.(BulkString)
		if !ok {
			return nil, false
		}
		bulkStrings = append(bulkStrings, element)
	}
	return bulkStrings, true
}
