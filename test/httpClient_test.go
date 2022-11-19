package test

import (
	"errors"
	"log"
	"my-cache2/cacheHttp"
	"my-cache2/mycache"
	"net/http"
	"testing"
)

func createGGroup() *mycache.Group {
	return mycache.NewGroup("group1", 2<<10, mycache.GetterFunc(func(key string) (mycache.RntValue, error) {
		s := db[key]
		if s == "" {
			return mycache.RntValue{}, errors.New("key is not find in db")
		}
		return mycache.RntValue{
			Bytes:  []byte(s),
			Expire: 2,
		}, nil
	}))
}

func Test_Server1(t *testing.T) {
	createGGroup()
	pool := cacheHttp.NewHttpServerPool("localhost:8801")
	log.Println("geecache is running at", "localhost:8801")
	http.ListenAndServe("localhost:8801", pool)
}
func Test_Server2(t *testing.T) {
	createGGroup()
	pool := cacheHttp.NewHttpServerPool("localhost:8802")
	log.Println("geecache is running at", "localhost:8802")
	http.ListenAndServe("localhost:8802", pool)

}
func Test_Server3(t *testing.T) {
	createGGroup()
	pool := cacheHttp.NewHttpServerPool("localhost:8803")
	log.Println("geecache is running at", "localhost:8803")
	http.ListenAndServe("localhost:8803", pool)
}

func Test_Client(t *testing.T) {
	gGroup := createGGroup()
	pool := cacheHttp.NewHttpClientPool("localhost:9999")
	gGroup.Register(pool)
	pool.Add("http://127.0.0.1:8801", "http://127.0.0.1:8802", "http://127.0.0.1:8803")
	cacheHttp.StartApiClient(gGroup, "localhost:9999")
}
