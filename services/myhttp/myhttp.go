package myhttp

import (
	"net/http"
	"time"
)

const timeout = time.Second * 60

var tr = &http.Transport{
	MaxIdleConns:      10,
	IdleConnTimeout:   30 * time.Second,
	DisableKeepAlives: true,
}

// Client 请求客户端
var Client = &http.Client{
	Transport: tr,
	Timeout:   timeout,
}
