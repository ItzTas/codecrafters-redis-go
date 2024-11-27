package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

type RDBKeyValue struct {
	key   string
	value string
	flag  byte
}

func printKeyValues(keyvals []RDBKeyValue) {
	fmt.Println("Key-Value len: ", len(keyvals))
	for _, kv := range keyvals {
		fmt.Printf("Key: %s, Value: %s, Flag: %d\n", kv.key, kv.value, kv.flag)
	}
}

func getValFromKeys(rdbs []RDBKeyValue, key string) (string, bool) {
	for _, kv := range rdbs {
		if kv.key == key {
			return kv.value, true
		}
	}
	return "", false
}

type RDBReader struct {
	file   *os.File
	noFile bool
}

func printHexdump(data []byte) {
	fmt.Println("Hexdump: ")
	for _, b := range data {
		if b >= 32 && b <= 126 {
			fmt.Printf("%c", b)
			continue
		}
		fmt.Print(".")
	}
	fmt.Println()
}

func extractReadable(data []byte) string {
	lines := bytes.Split(data, []byte("\n"))
	re := regexp.MustCompile(`\|(.+?)\|`)
	var result bytes.Buffer

	for _, line := range lines {
		matches := re.FindSubmatch(line)
		if len(matches) > 1 {
			result.Write(matches[1])
		}
	}

	return result.String()
}

func newRDBReader(path string) (*RDBReader, error) {
	noFile := false
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Could not create file: ", path)
		noFile = true
	}
	cmd := exec.Command("hexdump", "-C", path)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error in hexdump: " + err.Error())
	} else {
		fmt.Println("file created by codecrafters string: \n" + string(out))
		printHexdump(out)
		fmt.Println("readable: \n" + extractReadable(out))
		fmt.Println("")
	}
	return &RDBReader{file: file, noFile: noFile}, nil
}

func (r *RDBReader) readTillfe() error {
	buf := make([]byte, 1)
	for {
		_, err := r.file.Read(buf)
		if err != nil {
			return err
		}

		b := buf[0]
		if b == 0xfe {
			_, err := r.file.Read(buf)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

func (r *RDBReader) readDatabase() ([]RDBKeyValue, error) {
	if r.noFile {
		return nil, nil
	}
	keyValues := []RDBKeyValue{}

	err := r.readTillfe()
	if err != nil {
		return []RDBKeyValue{}, err
	}

	for {
		section, err := r.readDBSection()
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF reached")
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

	switch header[0] {
	case 0xfb:
		keyvals, err := r.readFB()
		if err != nil {
			fmt.Println("Error in reading FB: ", err)
		}
		return keyvals, err
	}

	fmt.Printf("Returning nil section: %x\n", header[0])
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

	hashSize := int(binary.LittleEndian.Uint16(hashSizeBytes))
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

	switch flagBytes[0] {
	case 0x00:
		return r.readString(hashSize)
	}

	return nil, nil
}

func (r *RDBReader) readString(sizeIndian int) ([]RDBKeyValue, error) {
	fmt.Printf("Reading string of size: %d\n", sizeIndian)

	var keyvals []RDBKeyValue

	for range sizeIndian {
		bufKeySize := make([]byte, 1)

		_, err := r.file.Read(bufKeySize)
		if err != nil {
			return []RDBKeyValue{}, err
		}

		keySize := int(bufKeySize[0])

		if keySize == 0 {
			_, err := r.file.Read(bufKeySize)
			if err != nil {
				return []RDBKeyValue{}, err
			}
			keySize = int(bufKeySize[0])
		}

		fmt.Println("keySize: ", keySize)

		buf := make([]byte, keySize)
		_, err = r.file.Read(buf)
		if err != nil {
			return []RDBKeyValue{}, err
		}

		key := string(buf)
		fmt.Println("Key: ", key)

		valueSizeBuf := make([]byte, 1)
		_, err = r.file.Read(valueSizeBuf)
		if err != nil {
			return []RDBKeyValue{}, err
		}

		valueSize := int(valueSizeBuf[0])
		fmt.Println("Value size: ", valueSize)
		bufValue := make([]byte, valueSize)
		_, err = r.file.Read(bufValue)
		if err != nil {
			return []RDBKeyValue{}, err
		}

		value := string(bufValue)
		fmt.Println("Value: ", value)

		keyValue := RDBKeyValue{
			key:   key,
			value: value,
			flag:  0x00,
		}
		keyvals = append(keyvals, keyValue)
	}

	printKeyValues(keyvals)
	return keyvals, nil
}
