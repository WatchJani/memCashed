package server

// header decoder
func Decode(payload []byte) (byte, uint32, uint32, uint32) {
	return payload[0], uint32(payload[1]), LittleEndian(payload[2:6]), LittleEndian(payload[6:10])
}

func LittleEndian(payload []byte) uint32 {
	return uint32(payload[0]) |
		uint32(payload[1])<<8 |
		uint32(payload[2])<<16 |
		uint32(payload[3])<<24
}
