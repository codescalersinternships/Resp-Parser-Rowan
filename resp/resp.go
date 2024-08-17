package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Type struct which represents type of resp being parsed
type Type byte

// const (
// 	SimpleString Type = '+'
// 	Error        Type = '-'
// 	Integer      Type = ':'
// 	BulkString   Type = '$'
// 	Array        Type = '*'
// )

// Value struct stores all info of resp parsed
type Value struct {
	typ     Type
	Integer int
	Str     []byte
	Array   []Value
	Err     error
	IsNull  bool
}

// RespReader struct is the buffer which is used to parse
type RespReader struct {
	reader *bufio.Reader
}

// NewReader returns a RespReader struct with the parsed reader
func NewReader(reader io.Reader) RespReader {
	return RespReader{reader: bufio.NewReader(reader)}
}

// ReadValue is the function which handles parsing the Value back
func (resp *RespReader) ReadValue() (Value, error) {
	char, err := resp.reader.ReadByte()
	if err != nil {
		return Value{IsNull: true}, err
	}
	fmt.Println(char == ':')
	if char == '*' {
		return resp.ReadValue()
	} else {
		switch char {
		case ':':
			return resp.readInteger()
		case '+':
			return resp.readString()
		case '$':
			return resp.readBulkString()
		case '-':
			return resp.readErrorMessage()
		}
	}
	return Value{IsNull: true}, fmt.Errorf("parsing error: beginning doesn't follow resp convensions")
}

// func (resp *RespReader) readArray() (Value, error) {

// }

func (resp *RespReader) readInteger() (Value, error) {
	fmt.Println("hello")
	num, err := resp.readInt()
	if err != nil {
		return Value{IsNull: true}, err
	}
	return Value{typ: ':', Integer: num}, err
}

func (resp *RespReader) readInt() (int, error) {
	line, err := resp.readLine()
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, err
	}
	return num, nil
}

// readLine converts a line ending with \r\n and returns the line without them, as an array of Bytes
func (resp *RespReader) readLine() (line []byte, err error) {
	for {
		bytes, err := resp.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		line = append(line, bytes...)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], err
}

// readBulkString parses bulk string specified with length
func (resp *RespReader) readBulkString() (Value, error) {
	l, err := resp.readInt()
	if err != nil {
		return Value{IsNull: true}, err
	}
	if l < 0 {
		return Value{typ: '$', IsNull: true}, fmt.Errorf("parsing error: string cannot have negative length")
	}
	// actual string length (added 2 bytes to read \r\n)
	strBytes := make([]byte, l+2)
	_, err = io.ReadFull(resp.reader, strBytes)
	if err != nil {
		return Value{IsNull: true}, err
	}
	if strBytes[l] != '\r' && strBytes[l+1] != '\n' {
		return Value{typ: '$', IsNull: true}, fmt.Errorf("parsing error: string doesn't end with the CRLF terminator")
	}
	return Value{typ: '$', Str: strBytes[:l]}, err
}

// readString parses simple string
func (resp *RespReader) readString() (Value, error) {
	line, err := resp.readLine()
	return Value{typ: '+', Str: line}, err
}

// readString parses error messages
func (resp *RespReader) readErrorMessage() (Value, error) {
	line, err := resp.readLine()
	return Value{typ: '-', Err: fmt.Errorf(string(line))}, err
}
