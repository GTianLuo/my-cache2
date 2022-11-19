package consistent

import (
	"fmt"
	"testing"
)

func TestConsistent(t *testing.T) {
	hashMap := NewConsistentMap(nil, 10)
	hashMap.Set("localhost:8081", "localhost:8082", "localhost:8083", "localhost:8084")
	fmt.Println(hashMap.Get("abc"))
	fmt.Println(hashMap.Get("abc"))
	fmt.Println(hashMap.Get("abc"))

	fmt.Println(hashMap.Get("efg"))
	fmt.Println(hashMap.Get("efg"))
	fmt.Println(hashMap.Get("efg"))

	fmt.Println(hashMap.Get("hio"))
	fmt.Println(hashMap.Get("hio"))
	fmt.Println(hashMap.Get("hio"))

	fmt.Println(hashMap.Get("klq"))
	fmt.Println(hashMap.Get("klq"))
	fmt.Println(hashMap.Get("klq"))
}
