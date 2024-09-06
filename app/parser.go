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
	invalidType = errors.New("Invalid input type")
	invalidResp = errors.New("Invalid resp format")
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
	r := Reader{
		r:     bytes.NewReader(payload),
		resps: []*RESP{},
	}

	return &r
}

type RESP struct {
	st    simbleType
	data  []byte
	count int
	array []*RESP
}

func (r *Reader) getCommand() string {
	return string(r.resps[0].data)
}

func (r *Reader) getArgs() []*RESP {
	return r.resps[1:]
}

func (r *Reader) readLine() ([]byte, error) {
	var bytes []byte
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			if err == io.EOF && len(bytes) > 0 {
				return bytes, nil
			}
			return []byte{}, err
		}

		if b == '\r' {
			b, err := r.r.ReadByte()
			if err != nil {
				return []byte{}, err
			}

			if b != '\n' {
				return []byte{}, invalidResp
			}

			return bytes, nil
		}

		bytes = append(bytes, b)
	}
}

func (r *Reader) readResp() error {
	typ, err := r.r.ReadByte()
	if err != nil {
		return err
	}
	switch simbleType(typ) {
	case BulkString:
		err := r.readBulk()
		return err
	case Array:
		err := r.readArray()
		return err
	case ErrorType:
		return nil
	case Integer:
		err := r.readInt()
		return err
	}

	return invalidType
}

func (r *Reader) readInt() error {
	bytes, err := r.readLine()
	if err != nil {
		return err
	}

	c, err := strconv.Atoi(string(bytes))
	if err != nil {
		return err
	}

	r.resps = append(r.resps, &RESP{
		st:    Integer,
		data:  bytes,
		count: c,
	})
	return nil
}

func (r *Reader) readArray() error {
	arLenByte, err := r.readLine()
	if err != nil {
		return err
	}

	arLen, err := strconv.Atoi(string(arLenByte))
	if err != nil {
		return err
	}

	if arLen == -1 {
		r.resps = append(r.resps, &RESP{
			st:    Array,
			data:  nil,
			count: arLen,
		})
		return nil
	}

	for range arLen {
		err := r.readResp()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Reader) readBulk() error {
	b, err := r.readLine()
	if err != nil {
		return err
	}

	bulkLen, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	if bulkLen == -1 {
		r.resps = append(r.resps, &RESP{
			st:    BulkString,
			data:  nil,
			count: bulkLen,
		})
		return nil
	}

	var bytes []byte

	for range bulkLen {
		b, err := r.r.ReadByte()
		if err != nil {
			return err
		}

		bytes = append(bytes, b)
	}

	crlf := make([]byte, 2)
	_, err = io.ReadFull(r.r, crlf)
	if err != nil {
		return err
	}

	if crlf[0] != '\r' || crlf[1] != '\n' {
		return fmt.Errorf("%v: %s", invalidResp, string(bytes))
	}

	r.resps = append(r.resps, &RESP{
		st:    BulkString,
		data:  bytes,
		count: bulkLen,
	})

	return nil
}
