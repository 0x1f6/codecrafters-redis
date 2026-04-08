package resp

import (
	"bufio"
	"fmt"
)

type SimpleError struct {
	msg string
}

func (e SimpleError) respValue() {}

func (e SimpleError) Error() string {
	return e.msg
}

func (e SimpleError) Serialize() []byte {
	buf := make([]byte, 0, len(e.msg)+3)
	buf = append(buf, '-')
	buf = append(buf, e.msg...)
	buf = append(buf, '\r', '\n')
	return buf
}

func NewSimpleError(msg string) SimpleError {
	return SimpleError{msg: msg}
}

func parseSimpleError(r *bufio.Reader) (SimpleError, error) {
	// Simple Error:
	// payload[0] == "-"
	// -<data>\r\n
	d, err := r.ReadBytes('\n')
	if err != nil {
		return SimpleError{}, err
	}

	if d[0] != '-' {
		return SimpleError{}, fmt.Errorf("Invalid prefix for RESP simple error: %c", d[0])
	}

	if len(d) <= 3 {
		return SimpleError{}, fmt.Errorf("Empty RESP simple error")
	}

	msg := d[1 : len(d)-2]
	return SimpleError{msg: string(msg)}, nil
}
