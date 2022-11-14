package http

import (
	"fmt"
	"log"
	"my-cache2/mycache"
	"net/http"
	"strings"
)

const defaultBasePath = "/my-cache/"

type HTTPServerPool struct {
	self     string
	basePath string
}

func NewHttpServerPool(self string) *HTTPServerPool {
	return &HTTPServerPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (h *HTTPServerPool) Log(logMessage string) {
	log.Printf("[Server %s]%s", h.self, logMessage)
}

func (h *HTTPServerPool) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.URL.Path, h.basePath) {
		h.Log(fmt.Sprintf("HTTPPool server is no expects url path:%s", request.URL.Path))
		http.Error(writer, fmt.Sprintf("HTTPPool server is no expects url path:%s", request.URL.Path), 404)
		return
	}
	split := strings.Split(request.URL.Path, "/")
	if len(split) < 4 {
		h.Log(fmt.Sprintf("HTTPPool server is no expects url path:%s", request.URL.Path))
		http.Error(writer, fmt.Sprintf("HTTPPool server is no expects url path:%s", request.URL.Path), 404)
		return
	}
	groupName := split[2]
	key := split[3]
	group, ok := mycache.GetGroup(groupName)
	if !ok {
		http.Error(writer, fmt.Sprintf("%s is not find in server", groupName), http.StatusBadRequest)
	}
	bytesValue, err := group.Get(key)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	writer.Write(bytesValue.ByteSlice())
}
