package test

import (
	"fmt"
	"testing"
)

type group struct {
	m map[string]string
}

func TestMap(t *testing.T) {
	//g := &group{}
	//fmt.Println(g.m)

	m := map[string]string{"1": "q", "2": "w", "3": "e", "4": "e", "5": "e"}
	for i, v := range m {
		fmt.Printf("%v  %v\n", i, v)
	}
}
