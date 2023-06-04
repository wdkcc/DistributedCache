package dCache

// ByteView 是按byte存储的数据
type ByteView struct {
	b []byte
}

// Len 数据的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回数据的副本作为切片
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 以字符串形式返回数据
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
