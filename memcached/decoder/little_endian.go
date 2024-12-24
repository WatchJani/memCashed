package decoder

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
