package parser

func DecodeLength(size []byte) int {
	return int(LittleEndianDecode(size))
}

// operation
// key length
// ttl - 4 byte
// payload - 4 byte
func Decode(payload []byte) (byte, uint32, uint32, uint32) {
	return payload[0], uint32(payload[1]), LittleEndianDecode(payload[2:6]), LittleEndianDecode(payload[6:10])
}
