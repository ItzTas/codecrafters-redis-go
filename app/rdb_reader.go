package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type RDBKeyValue struct {
	key   string
	value string
	flag  string
}

type RDBReader struct {
	file *os.File
}

func newRDBReader(path string) (*RDBReader, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Could not create file: ", path)
	}
	cmd := exec.Command("hexdump", "-C", path)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error in hexdump: " + err.Error())
	}
	fmt.Println("file created by codecrafters: " + string(out))
	return &RDBReader{file: file}, nil
}

func (r *RDBReader) readDatabase() ([]RDBKeyValue, error) {
	keyValues := []RDBKeyValue{}

	for {
		section, err := r.readDBSection()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading database section: ", err)
			return nil, err
		}
		if section == nil {
			continue
		}
		keyValues = append(keyValues, section...)
	}

	return keyValues, nil
}

func (r *RDBReader) readDBSection() ([]RDBKeyValue, error) {
	header := make([]byte, 1)
	_, err := r.file.Read(header)
	if err != nil {
		return nil, err
	}

	strHeader := fmt.Sprintf("%x", header[0])

	switch strHeader {
	case "fb":
		return r.readFB()
	}
	return nil, nil
}

func (r *RDBReader) readFB() ([]RDBKeyValue, error) {
	fmt.Println("Reading FB")
	hashSizeBytes := make([]byte, 2)
	_, err := r.file.Read(hashSizeBytes)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Hash size bytes: %x\n", hashSizeBytes)

	hashSize := int(hashSizeBytes[0])<<8 + int(hashSizeBytes[1])
	if hashSize <= 0 {
		return nil, errors.New("invalid FB size")
	}
	fmt.Println("Hash size: ", hashSize)

	flagBytes := make([]byte, 1)
	_, err = r.file.Read(flagBytes)
	if err != nil {
		return nil, err
	}
	flag := fmt.Sprintf("%x", flagBytes[0])
	fmt.Println("Flag: ", flag)

	switch flag {
	case "0":
		return r.readString(hashSize)
	}

	return []RDBKeyValue{}, nil
}

func (r *RDBReader) readString(size int) ([]RDBKeyValue, error) {
	fmt.Printf("Reading string of size: %d\n", size)
	data := make([]byte, size)
	_, err := r.file.Read(data)
	if err != nil {
		return nil, err
	}

	if len(data) < 2 {
		return nil, errors.New("invalid string data")
	}

	keySize := int(data[0])
	if len(data) < 1+keySize {
		return nil, errors.New("invalid key size")
	}

	key := string(data[1 : 1+keySize])
	valueSize := int(data[1+keySize])
	if len(data) < 2+keySize+valueSize {
		return nil, errors.New("invalid value size")
	}

	value := string(data[2+keySize : 2+keySize+valueSize])

	fmt.Println("Key:", key)
	fmt.Println("Value:", value)

	keyValue := RDBKeyValue{
		key:   key,
		value: value,
		flag:  "00",
	}

	return []RDBKeyValue{keyValue}, nil
}
