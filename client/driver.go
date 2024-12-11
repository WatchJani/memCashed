package client

import (
	"bytes"
	"encoding/binary"
)

// Operation - 1byte [x]
// Key length - 1byte [x]
// TTL - 4byte [x]
// body length - 4byte
// End of req - \r\n 2byte
func Req(key, value []byte, ttl int) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
	)

	//set operation
	if err = buf.WriteByte('S'); err != nil {
		// fmt.Println("Error writing operation:", err)
		return nil, err
	}

	//key length
	keyLength := uint8(len(key))
	err = binary.Write(&buf, binary.LittleEndian, keyLength)
	if err != nil {
		// fmt.Println("Error writing key length:", err)
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

	// fmt.Printf("Header: %v\n", buf.String())
	// fmt.Printf("Header length: %d\n", len(buf.Bytes()))
}
