package test

import (
	"errors"
	"fmt"
	"my-cache2/mycache"
	"testing"
	"time"
)

var db = map[string]string{
	"jok": "545",
	"klo": "323",
	"los": "232",
}

func TestGroup(t *testing.T) {
	mycache.NewGroup("group1", 2<<10, mycache.GetterFunc(func(key string) (mycache.RntValue, error) {
		s := db[key]
		if s == "" {
			return mycache.RntValue{}, errors.New("key is not find in db")
		}
		return mycache.RntValue{
			Bytes:  []byte(s),
			Expire: 2,
		}, nil
	}))

	group, _ := mycache.GetGroup("group1")
	fmt.Println(group.Get("jok"))
	fmt.Println(group.Get("klo"))
	fmt.Println(group.Get("jok"))
	fmt.Println(group.Get("klo"))
	time.Sleep(2 * time.Second)
	fmt.Println("sleep")
	fmt.Println(group.Get("jok"))
	fmt.Println(group.Get("klo"))
	fmt.Println(group.Get("jok"))
	fmt.Println(group.Get("klo"))
}
