package evn

import "encoding/binary"

func SimpleDecrypt(data []byte) []byte {
	magic := uint32(0xB36A210F)
	magicModifier := uint32(0xDEADBEEF)
	decryptedData := make([]byte, 0, len(data))

	// Process 4 bytes at a time
	for i := 0; i <= len(data)-4; i += 4 {
		// Convert 4 bytes to an integer
		chunk := binary.BigEndian.Uint32(data[i : i+4])
		// XOR with magic number
		j := chunk ^ magic
		// Convert back to bytes and append to decrypted_data
		decryptedChunk := make([]byte, 4)
		binary.BigEndian.PutUint32(decryptedChunk, j)
		decryptedData = append(decryptedData, decryptedChunk...)
		// Update magic number
		magic = (magic + magicModifier) & 0xFFFFFFFF
		magic ^= magicModifier
	}

	// Process any remaining bytes
	for i := len(data) - len(data)%4; i < len(data); i++ {
		decryptedData = append(decryptedData, data[i]^(byte(magic>>24)))
		magic = (magic << 8) & 0xFFFFFFFF
	}

	return decryptedData
}
