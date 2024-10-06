package main

import (
	"github.com/FDUTCH/obsidian/proxy/http_proxy"
	"log"
	"net/http"
	"time"
)

/*
	route-splitter example
*/

func main() {
	splitter := http_proxy.NewRouteSplitter(map[string]string{
		"/wiki": "wikipedia.org",
		"/dev":  "github.com",
		"/xxx":  "xxx.com",
	})
	handler := http_proxy.NewRouteSplitterHandler(splitter)
	http.ListenAndServe("localhost:8080", addMiddleware(handler, logger))
}

type writer struct {
	http.ResponseWriter
	handler func(status int)
}

func (w writer) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.handler(statusCode)
}

func logger(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	w = writer{w, func(status int) {
		log.Print(status, req.URL.Path, time.Since(now))
	}}
}

func addMiddleware(handler http.Handler, middleware ...func(http.ResponseWriter, *http.Request)) http.Handler {
	var fn http.HandlerFunc = handler.ServeHTTP

	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			for _, handlerFunc := range middleware {
				handlerFunc(writer, request)
			}
			fn(writer, request)
		},
	)
}
