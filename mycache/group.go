package mycache

import (
	"log"
	my_cache2 "my-cache2"
	"my-cache2/nodes"
	"sync"
)

type getter interface {
	get(key string) (RntValue, error)
}

type GetterFunc func(key string) (RntValue, error)

func (g GetterFunc) get(key string) (RntValue, error) {
	return g(key)
}

type RntValue struct {
	Bytes  []byte
	Expire int64
}

type Group struct {
	name       string
	cache      *mainCache
	getter     getter
	nodePicker nodes.NodePicker
}

var (
	rw     sync.RWMutex
	groups map[string]*Group = make(map[string]*Group)
)

func NewGroup(name string, maxBytes uint64, getter getter) *Group {
	rw.Lock()
	defer rw.Unlock()
	if getter == nil {
		panic("the getter is not allowed to be nil")
	}
	group := &Group{
		name:   name,
		getter: getter,
		cache:  &mainCache{maxBytes: maxBytes},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) (*Group, bool) {
	rw.RLock()
	defer rw.RUnlock()
	group := groups[name]
	if group != nil {
		return group, true
	} else {
		return nil, false
	}
}

func (g *Group) Register(picker nodes.NodePicker) {
	if g.nodePicker != nil {
		panic("The number of registrations cannot exceed once")
	}
	g.nodePicker = picker
}

func (g *Group) Get(key string) (my_cache2.BytesValue, error) {
	bytesValue, ok := g.cache.get(key)
	if ok {
		log.Printf("[MyCache] %s is hit in cache\n", key)
		return bytesValue, nil
	}
	return g.Load(key)
}

func (g *Group) Load(key string) (value my_cache2.BytesValue, err error) {
	if g.nodePicker != nil {
		if node, ok := g.nodePicker.PickNode(key); ok {
			if value, err = g.getFromPeer(node, key); err == nil {
				return value, nil
			}
			log.Println("failed to load value form node ", err)
		}
	}
	return g.loadLocally(key)
}

func (g *Group) getFromPeer(node nodes.NodeGetter, key string) (bytesValue my_cache2.BytesValue, err error) {
	if bytes, err := node.Get(g.name, key); err != nil {
		return my_cache2.BytesValue{Bytes: bytes}, nil
	}
	return bytesValue, nil
}

func (g *Group) loadLocally(key string) (my_cache2.BytesValue, error) {
	rntValue, err := g.getter.get(key)
	if err != nil {
		return my_cache2.BytesValue{}, err
	}
	log.Printf("[Slow DB] %s is searched in DB", key)
	value := my_cache2.BytesValue{Bytes: rntValue.Bytes}
	err = g.syncToCache(key, value, rntValue.Expire)
	return value, err
}

func (g *Group) syncToCache(key string, value my_cache2.BytesValue, expire int64) error {
	return g.cache.add(key, value, expire)
}
