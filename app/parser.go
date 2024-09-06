package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type simbleType rune

var (
	invalidType = errors.New("invalid input type")
	invalidResp = errors.New("invalid resp format")
)

const (
	BulkString   simbleType = '$'
	SimpleString simbleType = '+'
	Integer      simbleType = ':'
	Array        simbleType = '*'
	ErrorType    simbleType = '-'
)

type Reader struct {
	r     *bytes.Reader
	resps []*RESP
}

func NewReader(payload []byte) *Reader {
	return &Reader{
		r:     bytes.NewReader(payload),
		resps: []*RESP{},
	}
}

type RESP struct {
	st    simbleType
	data  []byte
	count int
	array []*RESP
}

func (r *Reader) getCommand() string {
	if len(r.resps[0].array) == 0 {
		return ""
	}

	return string(r.resps[0].array[0].data)
}

func (r *Reader) PrintResps() {
	for i, resp := range r.resps {
		fmt.Printf("RESP #%d:\n", i+1)
		fmt.Printf("  Type: %c\n", resp.st)
		fmt.Printf("  Data: %s\n", string(resp.data))

		if resp.st == Array {
			fmt.Printf("  Array elements: %d\n", resp.count)
			for j, item := range resp.array {
				fmt.Printf("    Element #%d: %v\n", j+1, *item)
			}
		} else {
			fmt.Printf("  Count: %d\n", resp.count)
		}

		fmt.Printf("  Full RESP object: %+v\n", resp)
	}
}

func (r *Reader) getArgs() []*RESP {
	return r.resps[0].array[1:]
}

func (r *Reader) readLine() ([]byte, error) {
	var bytes []byte
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF && len(bytes) > 0 {
				return bytes, nil
			}
			return nil, err
		}

		if b == '\r' {
			b, err := r.r.ReadByte()
			if err != nil {
				return nil, err
			}

			if b != '\n' {
				return nil, invalidResp
			}

			return bytes, nil
		}

		bytes = append(bytes, b)
	}
}

func (r *Reader) readResp() (*RESP, error) {
	typ, err := r.r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch simbleType(typ) {
	case BulkString:
		return r.readBulk()
	case Array:
		return r.readArray()
	case ErrorType:
		return r.readError()
	case Integer:
		return r.readInt()
	case SimpleString:
		return r.readSimpleString()
	default:
		return nil, invalidType
	}
}

func (r *Reader) readInt() (*RESP, error) {
	bytes, err := r.readLine()
	if err != nil {
		return nil, err
	}

	value, err := strconv.Atoi(string(bytes))
	if err != nil {
		return nil, err
	}

	resp := &RESP{
		st:    Integer,
		data:  bytes,
		count: value,
	}
	return resp, nil
}

func (r *Reader) readArray() (*RESP, error) {
	arLenByte, err := r.readLine()
	if err != nil {
		return nil, err
	}

	arLen, err := strconv.Atoi(string(arLenByte))
	if err != nil {
		return nil, err
	}

	resp := &RESP{
		st:    Array,
		count: arLen,
		array: make([]*RESP, arLen),
	}

	if arLen == -1 {
		r.resps = append(r.resps, resp)
		return resp, nil
	}

	for i := 0; i < arLen; i++ {
		itemResp, err := r.readResp()
		if err != nil {
			return nil, err
		}
		resp.array[i] = itemResp
	}

	r.resps = append(r.resps, resp)
	return resp, nil
}

func (r *Reader) readBulk() (*RESP, error) {
	b, err := r.readLine()
	if err != nil {
		return nil, err
	}

	bulkLen, err := strconv.Atoi(string(b))
	if err != nil {
		return nil, err
	}

	resp := &RESP{
		st:    BulkString,
		count: bulkLen,
	}

	if bulkLen == -1 {
		r.resps = append(r.resps, resp)
		return resp, nil
	}

	var bytes []byte
	for i := 0; i < bulkLen; i++ {
		b, err := r.r.ReadByte()
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
	}

	crlf := make([]byte, 2)
	_, err = io.ReadFull(r.r, crlf)
	if err != nil {
		return nil, err
	}

	if crlf[0] != '\r' || crlf[1] != '\n' {
		return nil, fmt.Errorf("%v: %s", invalidResp, string(bytes))
	}

	resp.data = bytes
	return resp, nil
}

func (r *Reader) readSimpleString() (*RESP, error) {
	bytes, err := r.readLine()
	if err != nil {
		return nil, err
	}

	resp := &RESP{
		st:    SimpleString,
		data:  bytes,
		count: len(bytes),
	}
	r.resps = append(r.resps, resp)
	return resp, nil
}

func (r *Reader) readError() (*RESP, error) {
	bytes, err := r.readLine()
	if err != nil {
		return nil, err
	}

	resp := &RESP{
		st:    ErrorType,
		data:  bytes,
		count: len(bytes),
	}
	r.resps = append(r.resps, resp)
	return resp, nil
}
