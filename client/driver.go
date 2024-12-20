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

// var bodyLength uint32 = 4

// 	bloc := []byte{
// 		byte(bodyLength & 0xFF),
// 		byte((bodyLength >> 8) & 0xFF),
// 		byte((bodyLength >> 16) & 0xFF),
// 		byte((bodyLength >> 24) & 0xFF),
// 	}

func Encode(operation byte, key, value []byte, ttl int) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
	)

	payloadSize := uint32(len(value) + len(key) + 10)
	err = binary.Write(&buf, binary.LittleEndian, payloadSize)
	if err != nil {
		return nil, err
	}

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

// operation
// key length
// ttl - 4 byte
// payload - 4 byte
func Decode(payload []byte) (byte, uint32, uint32, uint32) {
	return payload[0], uint32(payload[1]), LittleEndian(payload[2:6]), LittleEndian(payload[6:10])
}

func TestDecode(payload []byte) (uint32, byte, uint32, uint32, uint32) {
	return LittleEndian(payload[:4]), payload[4], uint32(payload[5]), LittleEndian(payload[6:10]), LittleEndian(payload[10:14])
}

func LittleEndian(payload []byte) uint32 {
	return uint32(payload[0]) |
		uint32(payload[1])<<8 |
		uint32(payload[2])<<16 |
		uint32(payload[3])<<24
}

func DecodeLength(size []byte) int {
	return int(LittleEndian(size))
}
