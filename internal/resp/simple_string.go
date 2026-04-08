package resp

import (
	"bufio"
	"fmt"
)

type SimpleString struct {
	msg string
}

func (s SimpleString) respValue() {}

func (s SimpleString) String() string {
	return s.msg
}

func (s SimpleString) Serialize() []byte {
	buf := make([]byte, 0, len(s.msg)+3)
	buf = append(buf, '+')
	buf = append(buf, s.msg...)
	buf = append(buf, '\r', '\n')
	return buf
}

func NewSimpleString(msg string) SimpleString {
	return SimpleString{msg: msg}
}

func parseSimpleString(r *bufio.Reader) (SimpleString, error) {
	// Simple Strings:
	// payload[0] == "+"
	// +<data>\r\n
	d, err := r.ReadBytes('\n')
	if err != nil {
		return SimpleString{}, err
	}

	if d[0] != '+' {
		return SimpleString{}, fmt.Errorf("Invalid prefix for RESP simple string: %c", d[0])
	}

	msg := d[1 : len(d)-2]
	return SimpleString{msg: string(msg)}, nil
}
