package http

import (
	"fmt"
	"io"
	"my-cache2/consistent"
	"my-cache2/nodes"
	"net/http"
	"sync"
)

var (
	defaultReplicate = 50
)

type httpGetter struct {
	baseUrl string
}

type HttpClientPool struct {
	*HTTPServerPool
	mu          sync.Mutex
	nodes       *consistent.ConsistentMap
	httpGetters map[string]*httpGetter
}

func NewHttpClientPool(self string) *HttpClientPool {
	return &HttpClientPool{
		HTTPServerPool: NewHttpServerPool(self),
		nodes:          consistent.NewConsistentMap(nil, defaultReplicate),
	}
}

func (h *HttpClientPool) Add(nodes ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.nodes.Set(nodes...)
	if h.httpGetters == nil {
		h.httpGetters = make(map[string]*httpGetter)
	}
	for _, node := range nodes {
		h.httpGetters[node] = &httpGetter{
			baseUrl: node + h.basePath,
		}
	}
}

func (h *HttpClientPool) PickNode(key string) (nodes.NodeGetter, bool) {
	if node := h.nodes.Get(key); node != "" && node != h.self {
		return h.httpGetters[node], true
	}
	return nil, false
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	url := fmt.Sprintf("%v%v/%v", h.baseUrl, group, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Status != "200 ok" {
		return nil, fmt.Errorf("Wrong return status code:" + resp.Status)
	}
	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
