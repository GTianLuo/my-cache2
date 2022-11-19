package consistent

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type hash func([]byte) uint32

type ConsistentMap struct {
	hash      hash
	replicate int
	keys      []int
	m         map[int]string
}

func NewConsistentMap(fn hash, replicate int) *ConsistentMap {
	if replicate < 0 {
		panic("replicate must must be greater than 0")
	}
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &ConsistentMap{
		replicate: replicate,
		hash:      fn,
	}
}

func (cMap *ConsistentMap) Set(nodes ...string) {
	if cMap.m == nil {
		cMap.m = make(map[int]string, len(nodes))
	}
	for _, node := range nodes {
		for i := 0; i < cMap.replicate; i++ {
			hash := int(cMap.hash([]byte(strconv.Itoa(i) + node)))
			cMap.keys = append(cMap.keys, hash)
			cMap.m[hash] = node
		}

	}
	sort.Ints(cMap.keys)
}

func (cMap *ConsistentMap) Get(key string) string {
	hash := int(cMap.hash([]byte(key)))
	i := sort.Search(len(cMap.keys), func(i int) bool {
		return cMap.keys[i] >= hash
	})
	return cMap.m[cMap.keys[i%len(cMap.keys)]]
}
