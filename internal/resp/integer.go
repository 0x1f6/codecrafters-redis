package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

type Integer struct {
	value int
}

func (i Integer) respValue() {}

func (i Integer) String() string {
	return fmt.Sprint(i.value)
}

func (i Integer) Serialize() []byte {
	valueStr := strconv.Itoa(i.value)
	buf := make([]byte, 0, len(valueStr)+3)
	buf = append(buf, ':')
	buf = append(buf, valueStr...)
	buf = append(buf, '\r', '\n')
	return buf
}

func NewInteger(value int) Integer {
	return Integer{value: value}
}

func parseInteger(r *bufio.Reader) (Integer, error) {
	// Integer:
	// payload[0] == ":"
	// :[<+|->]<value>\r\n
	d, err := r.ReadBytes('\n')
	if err != nil {
		return Integer{}, err
	}

	if d[0] != ':' {
		return Integer{}, fmt.Errorf("Invalid prefix for RESP integer: %c", d[0])
	}

	var valueStr string
	if d[1] == '-' {
		valueStr = string(d[2 : len(d)-2])
	} else {
		valueStr = string(d[1 : len(d)-2])
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return Integer{}, fmt.Errorf("could not parse value as integer")
	}
	return Integer{value: value}, nil
}
