package client

var EmptyByte []byte = []byte{}

// Operation - 1byte [x]
// Key length - 1byte [x]
// TTL - 4byte [x]
// body length - 4byte
// End of req - \r\n 2byte
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

func LittleEndianEncode(payload []byte, num uint32) int {
	payload[0] = byte(num & 0xFF)
	payload[1] = byte((num >> 8) & 0xFF)
	payload[2] = byte((num >> 16) & 0xFF)
	payload[3] = byte((num >> 24) & 0xFF)

	return 4
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

// operation
// key length
// ttl - 4 byte
// payload - 4 byte
func Decode(payload []byte) (byte, uint32, uint32, uint32) {
	return payload[0], uint32(payload[1]), LittleEndianDecode(payload[2:6]), LittleEndianDecode(payload[6:10])
}

func LittleEndianDecode(payload []byte) uint32 {
	return uint32(payload[0]) |
		uint32(payload[1])<<8 |
		uint32(payload[2])<<16 |
		uint32(payload[3])<<24
}

func DecodeLength(size []byte) int {
	return int(LittleEndianDecode(size))
}
