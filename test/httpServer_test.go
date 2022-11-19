package test

import (
	"errors"
	"log"
	http2 "my-cache2/cacheHttp"
	"my-cache2/mycache"
	"net/http"
	"testing"
)

func TestHttpServer(t *testing.T) {
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
	serverPool := http2.NewHttpServerPool("localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", serverPool))
}
