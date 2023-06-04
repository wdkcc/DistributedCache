package LRU

import "container/list"

// Cache 是LRU缓存
type Cache struct {
	maxBytes  int64                         // 最大内存容量
	nbytes    int64                         // 当前存储的数据量
	ll        *list.List                    // 存储数据的链表
	cache     map[string]*list.Element      // 字典
	OnEvicted func(key string, value Value) // 可选项，回调函数
}

// entry 是存储的数据
type entry struct {
	key   string
	value Value
}

// Value 存储的值必须要实现Len方法计算其占多少字节
type Value interface {
	Len() int
}

// New 创建缓存
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 根据Key值查找数据
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 淘汰缓存
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增数据
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 待添加数据已在缓存中，则将数据移至队头并更新
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 不在缓存中，插入链表并建立map映射
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 添加数据后，达到最大限制，就淘汰旧数据
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 缓存数据长度
func (c *Cache) Len() int {
	return c.ll.Len()
}
