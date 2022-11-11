package test

import (
	"fmt"
	"my-cache2/evict"
	"testing"
	"time"
)

type v struct {
	s string
}

func (v v) Len() int {
	return len(v.s)
}

// 测试对url剔除、对设置超时时间和未设置超时时间进行分组
func TestCache1(t *testing.T) {
	cache := evict.NewCache(20)
	cache.Add("12", v{"ab"}, 2)
	cache.Add("34", v{"ab"}, 2)
	cache.Add("56", v{"ab"}, 2)
	cache.Add("78", v{"ab"}, 2)
	cache.Add("90", v{"ab"}, 2)
	cache.Add("91", v{"ab"}, 2)
	cache.Add("92", v{"ab"}, 2)
	cache.Add("93", v{"ab"}, -1)
	cache.Add("94", v{"ab"}, -1)
	cache.Add("95", v{"ab"}, -1)
	cache.Print()
}

// 测式定时随机剔除
func Test(t *testing.T) {
	cache := evict.NewCache(20)
	cache.DeleteExpired()
	cache.Add("12", v{"ab"}, 2)
	cache.Add("34", v{"ab"}, 2)
	cache.Add("56", v{"ab"}, 2)
	cache.Add("78", v{"ab"}, 2)
	cache.Add("90", v{"ab"}, 2)
	cache.Add("95", v{"ab"}, -1)
	cache.Print()
	time.Sleep(4 * time.Second)
	cache.Print()
}

func TestCache3(t *testing.T) {
	cache := evict.NewCache(20)
	cache.Add("12", v{"ab"}, 2)
	cache.Add("34", v{"ab"}, 2)
	cache.Add("56", v{"ab"}, 2)
	cache.Add("78", v{"ab"}, 2)
	cache.Add("90", v{"ab"}, 2)
	fmt.Println(cache.Get("12"))
	fmt.Println(cache.Get("34"))
	fmt.Println(cache.Get("56"))
	fmt.Println(cache.Get("78"))
	fmt.Println(cache.Get("90"))
	time.Sleep(4 * time.Second)
	fmt.Println(cache.Get("12"))
	fmt.Println(cache.Get("34"))
	fmt.Println(cache.Get("56"))
	fmt.Println(cache.Get("78"))
	fmt.Println(cache.Get("90"))
	cache.Print()
}
