package endian


func Flip(bytesToFlip []byte) []byte {
	allFlippedBytes := make([]byte, 0, len(bytesToFlip))
	for i := 0; i < len(bytesToFlip); i += 2 {
		if i+2 <= len(bytesToFlip) {
			allFlippedBytes = append(allFlippedBytes, flipTwoBytes(bytesToFlip[i:i+2])...)
		}
	}
	return allFlippedBytes
}

func flipTwoBytes(b []byte) []byte {
	if len(b) != 2 {
		return b
	}
	return []byte{b[1], b[0]}
}
