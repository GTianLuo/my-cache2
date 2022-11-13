package my_cache2

type BytesValue struct {
	Bytes []byte
}

func (v BytesValue) Len() int {
	return len(v.Bytes)
}

func (v BytesValue) String() string {
	return string(v.Bytes)
}

func (v BytesValue) ByteSlice() []byte {
	return cloneBytes(v.Bytes)
}

func cloneBytes(bytes []byte) []byte {
	b := make([]byte, len(bytes))
	copy(b, bytes)
	return b
}
