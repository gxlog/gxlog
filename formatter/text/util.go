package text

func cloneBytes(slice []byte) []byte {
	clone := make([]byte, len(slice))
	copy(clone, slice)
	return clone
}
