package api

import (
	"main/model"
	"net/http"
)

// https://api.cloudflare.com/#dns-records-for-a-zone-list-dns-records
const (
	// BASEURL cloudflare的基础请求地址
	BASEURL = "https://api.cloudflare.com/client/v4"
)

// 获取认证字段
func getAuthHeader(config *model.Config) (header http.Header) {
	header = http.Header{}
	header.Add("X-Auth-Email", config.AuthEmail)
	header.Add("X-Auth-Key", config.AuthKey)
	header.Add("Content-type", "application/json")

	return header
}
