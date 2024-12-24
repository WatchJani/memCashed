package decoder

func LittleEndianDecode(payload []byte) uint32 {
	return uint32(payload[0]) |
		uint32(payload[1])<<8 |
		uint32(payload[2])<<16 |
		uint32(payload[3])<<24
}

func LittleEndianEncode(payload []byte, num uint32) int {
	payload[0] = byte(num & 0xFF)
	payload[1] = byte((num >> 8) & 0xFF)
	payload[2] = byte((num >> 16) & 0xFF)
	payload[3] = byte((num >> 24) & 0xFF)

	return 4
}
