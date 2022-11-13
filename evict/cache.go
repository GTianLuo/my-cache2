package evict

import (
	"container/list"
	"fmt"
	"time"
)

type Value interface {
	Len() int
}

type entity struct {
	key string
	v   Value
	ddl int64
}

type Cache struct {
	maxBytes      uint64
	uBytes        uint64
	ll            *list.List
	allDataMap    map[string]*list.Element
	expireDataMap map[string]*list.Element
}

func NewCache(maxBytes uint64) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		ll:         new(list.List),
		allDataMap: make(map[string]*list.Element),
	}
}

func (c *Cache) RemoveOldest() {
	back := c.ll.Back()
	c.ll.Remove(back)
	e := back.Value.(*entity)
	delete(c.allDataMap, e.key)
	delete(c.expireDataMap, e.key)
	c.uBytes -= uint64(back.Value.(*entity).v.Len()) + uint64(len(e.key))
}

func (c *Cache) Add(key string, v Value, expire int64) error {
	if c.allDataMap[key] != nil {
		c.remove(key)
	}
	vBytes := uint64(v.Len()) + uint64(len([]byte(key)))
	if vBytes > c.maxBytes-c.uBytes {
		if vBytes > c.maxBytes {
			return fmt.Errorf("%s is not find in cache", key)
		}
		c.RemoveOldest()
		return c.Add(key, v, expire)
	}
	var ddl int64 = -1
	if expire > 0 {
		ddl = time.Now().Unix() + expire
	}
	e := c.ll.PushFront(&entity{key, v, ddl})
	c.uBytes += vBytes
	c.allDataMap[key] = e
	if expire > 0 {
		if c.expireDataMap == nil {
			c.expireDataMap = make(map[string]*list.Element)
		}
		c.expireDataMap[key] = e
	}
	return nil
}

func (c *Cache) remove(key string) {
	element := c.allDataMap[key]
	c.ll.MoveToBack(element)
	c.RemoveOldest()
}

func (c *Cache) Get(key string) (Value, bool) {
	element := c.allDataMap[key]
	if element == nil {
		return nil, false
	}
	entity := element.Value.(*entity)
	if entity.ddl > 0 && entity.ddl < time.Now().Unix() {
		c.remove(key)
		return nil, false
	}
	c.ll.MoveToFront(element)
	return entity.v, true
}

func (c *Cache) DeleteExpired() {
	go func() {
		for true {
			if c.expireDataMap == nil {
				time.Sleep(1 * time.Second)
				continue
			}
			count := 20
			expired := 0
			for _, v := range c.expireDataMap {
				if count <= 0 {
					break
				}
				e := v.Value.(*entity)
				if e.ddl <= time.Now().Unix() {
					expired++
					c.remove(e.key)
				}
				count--
			}
			if expired < 5 {
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (c *Cache) Print() {
	fmt.Println("allDataMap:")
	for _, v := range c.allDataMap {
		fmt.Printf("%v  ", v.Value.(*entity).key)
	}
	fmt.Println("\nexpireDataMap")
	for _, v := range c.expireDataMap {
		fmt.Printf("%v  ", v.Value.(*entity).key)
	}
	fmt.Println()
}

func (c *Cache) Len() int {
	return len(c.allDataMap)
}
