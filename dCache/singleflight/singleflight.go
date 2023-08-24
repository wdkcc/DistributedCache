package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex // protects m
	m  map[string]*call
}

// Do 让同一小段时间访问相同key的请求，只访问一次的缓存和数据库
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	// 1. 第一次进来时对g加锁，其它协程会被阻塞在这里，直至锁释放
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 2. 第一个进来的协程会发现key值不存在，后面进来的协程运行到这里时key值已经存在，就直接取出来
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	// 3. 对c加锁，直至读取完数据存入c中后再解锁
	c.wg.Add(1)
	g.m[key] = c
	// 4. 解锁g，注意此时会将后续操作相同key的协程全部放入第二步的if语句中，等待函数执行完数据写入c
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key) // 从map中删除该key对应的映射，因为运行到这的时候，同一时间访问该key的协程在第二步中都已得到数据
	g.mu.Unlock()

	return c.val, c.err
}
