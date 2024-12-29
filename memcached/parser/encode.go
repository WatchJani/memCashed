package parser

var EmptyByte []byte = []byte{}

func Set(key, value []byte, ttl int) ([]byte, error) {
	return Encode('S', key, value, ttl)
}

func SetDriver(key, value, reqID []byte, ttl int) ([]byte, error) {
	return Encode('S', key, value, ttl)
}

func Get(key []byte) ([]byte, error) {
	return Encode('G', key, EmptyByte, 0)
}

func Delete(key []byte) ([]byte, error) {
	return Encode('D', key, EmptyByte, 0)
}

func Encode(operation byte, key, value []byte, ttl int) ([]byte, error) {
	payloadSize := uint32(len(value) + len(key) + 10)

	buf := make([]byte, payloadSize+4)
	offset := 0

	// payloadSize := uint32(len(value) + len(key) + 10)
	offset += LittleEndianEncode(buf[offset:offset+4], payloadSize)

	buf[offset] = operation
	offset++

	//key length
	buf[offset] = uint8(len(key))
	offset++

	//set ttl
	TTL := uint32(ttl)
	offset += LittleEndianEncode(buf[offset:offset+4], TTL)

	//set body length
	bodyLength := uint32(len(value))
	offset += LittleEndianEncode(buf[offset:offset+4], bodyLength)

	offset += copy(buf[offset:], key)

	offset += copy(buf[offset:], value)

	return buf, nil
}
