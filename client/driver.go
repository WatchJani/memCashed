package client

import (
	"bytes"
	"encoding/binary"
)

var EmptyByte []byte = []byte{}

// Operation - 1byte [x]
// Key length - 1byte [x]
// TTL - 4byte [x]
// body length - 4byte
// End of req - \r\n 2byte
func Set(key, value []byte, ttl int) ([]byte, error) {
	return Encode('S', key, value, ttl)
}

func Get(key []byte) ([]byte, error) {
	return Encode('G', key, EmptyByte, 0)
}

func Delete(key []byte) ([]byte, error) {
	return Encode('D', key, EmptyByte, 0)
}

func Encode(operation byte, key, value []byte, ttl int) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
	)

	//set operation
	if err = buf.WriteByte(operation); err != nil {
		return nil, err
	}

	//key length
	keyLength := uint8(len(key))
	err = binary.Write(&buf, binary.LittleEndian, keyLength)
	if err != nil {
		return nil, err
	}

	//set ttl
	TTL := uint32(ttl)
	err = binary.Write(&buf, binary.LittleEndian, TTL)
	if err != nil {
		return nil, err
	}

	//set body length
	bodyLength := uint32(len(value))
	err = binary.Write(&buf, binary.LittleEndian, bodyLength)
	if err != nil {
		return nil, err
	}

	//set real key
	if _, err = buf.Write(key); err != nil {
		return nil, err
	}

	//Set real body
	if _, err = buf.Write(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
