package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps 字节映射为uint32类型
type Hash func(data []byte) uint32

// Map 一致性hash核心结构
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 倍数
	keys     []int          // 排序的虚拟节点keys值集合
	hashMap  map[int]string // 虚拟节点和真实节点的映射表
}

// New 实例化 Map
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加键
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 获取最近的节点作为缓存的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 二分查找最近的key索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
