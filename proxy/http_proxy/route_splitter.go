package http_proxy

import (
	"net/http"
	"strings"
)

type routeSplitter map[string]string

func (r routeSplitter) Address(route string) string {
	return r[route]
}

// RouteSplitter splits http requests by there routes between hosts
type RouteSplitter interface {
	Address(route string) string
}

// RouteSplitterHandler  http.Handler implementation for splitting http request between hosts
type RouteSplitterHandler struct {
	splitter RouteSplitter
	c        *cache
}

func NewRouteSplitterHandler(splitter RouteSplitter) *RouteSplitterHandler {
	return &RouteSplitterHandler{splitter: splitter, c: new(cache)}
}

func (r RouteSplitterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := strings.Split(request.URL.Path, "/")
	var addr string
	if request.URL.Path == "" {
		addr = r.splitter.Address("/")
	} else {
		request.URL.Path = strings.TrimPrefix(request.URL.Path, "/"+path[1])
		addr = r.splitter.Address("/" + path[1])
	}

	proxy, err := r.c.getProxy(addr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	proxy.ServeHTTP(writer, request)
}

func NewRouteSplitter(route map[string]string) RouteSplitter {
	return routeSplitter(route)
}
