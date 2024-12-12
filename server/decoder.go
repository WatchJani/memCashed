package server

func Decode(payload []byte) (byte, uint32, uint32, uint32) {
	operation := payload[0]

	keyLength := uint32(payload[1])

	ttl := uint32(payload[2]) |
		uint32(payload[3])<<8 |
		uint32(payload[4])<<16 |
		uint32(payload[5])<<24

	bodyLength := uint32(payload[6]) |
		uint32(payload[7])<<8 |
		uint32(payload[8])<<16 |
		uint32(payload[9])<<24

	return operation, keyLength, ttl, bodyLength
}
