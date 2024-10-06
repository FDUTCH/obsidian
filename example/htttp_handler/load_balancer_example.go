package main

import (
	"github.com/FDUTCH/obsidian/proxy/balance"
	"github.com/FDUTCH/obsidian/proxy/http_proxy"
	"log"
	"net/http"
)

/*
	Simple load balance proxy for http
*/

func main() {
	loadBalancer := balance.NewSimpleLoadBalancer("www.wikipedia.org", "bedrock.dev", "pornhub.com", "example.com")
	handler := http_proxy.NewBalanceHandler(loadBalancer)
	log.Fatal(http.ListenAndServe("localhost:8080", handler))
}
