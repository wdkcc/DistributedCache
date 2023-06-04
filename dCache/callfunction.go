package dCache

// Getter 接收回调函数的接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 回调函数的类型，输入key，得到数据
type GetterFunc func(key string) ([]byte, error)

// Get 接口方法实现
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
