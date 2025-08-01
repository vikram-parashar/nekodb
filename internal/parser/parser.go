package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	BULK    = '$'
	ARRAY   = '*'
)

type RespParser struct {
	r *bufio.Reader
}

func NewParser(rd io.Reader) *RespParser {
	return &RespParser{
		r: bufio.NewReader(rd),
	}
}

type DataType struct {
	Name    string
	Bulk     string
	Arr     []DataType
}

func (rp *RespParser) Read() (DataType, error) {
	b, err := rp.r.ReadByte()
	if err != nil {
		return DataType{}, err
	}

	switch b {
	case BULK:
		return rp.ReadBulk()
	case ARRAY:
		return rp.ReadArray()
	default:
		fmt.Println(string(b), b)
		return DataType{}, fmt.Errorf("Invalid DataType")
	}
}

func (rp *RespParser) ReadArray() (DataType, error) {
	n, err := rp.ReadInt()
	if err != nil {
		return DataType{}, nil
	}

	arr := make([]DataType, n)
	for i := range arr {
		arr[i], err = rp.Read()
		if err != nil {
			return DataType{}, err
		}
	}
	return DataType{
		Name: "array",
		Arr:  arr,
	}, err
}

func (rp *RespParser) ReadBulk() (DataType, error) {
	n, err := rp.ReadInt()
	if err != nil {
		return DataType{}, nil
	}

	bulk := make([]byte, n)
	_, err = rp.r.Read(bulk)
	if err != nil {
		return DataType{}, err
	}

	//skip crlf
	rp.ReadLine()

	return DataType{
		Name: "bulk",
		Bulk: string(bulk),
	}, nil

}

func (r *RespParser) ReadInt() (n int, err error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(line))
}

func (r *RespParser) ReadLine() (line []byte, err error) {
	line, err = r.r.ReadBytes('\n')

	//if err or not remove \r\n
	if err != nil || len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
		return nil, err
	}

	line = line[:len(line)-2]
	return line, err
}
