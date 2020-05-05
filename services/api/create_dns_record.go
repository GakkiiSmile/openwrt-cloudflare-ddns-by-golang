package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/services/myhttp"
	"net/http"
	"net/url"
)

// CreateDNSRecordResponse 响应结果
type createDNSRecordResponse struct {
	Success  bool            `json:"success"`
	Errors   []interface{}   `json:"errors"`
	Messages []interface{}   `json:"messages"`
	Result   model.DNSRecord `json:"result"`
}

type postCreateDNSRecordJSONData struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Proxied  bool   `json:"proxied"`
}

// CreateDNSRecord 创建
func CreateDNSRecord(config *model.Config, zoneID string, Type string, Name string, Content string, Proxied bool) (result model.DNSRecord, err error) {
	url, err := url.Parse(BASEURL + "/zones/" + zoneID + "/dns_records")

	result = model.DNSRecord{}

	if err != nil {
		return result, err
	}

	jsonData := &postCreateDNSRecordJSONData{
		Type:     Type,
		Name:     Name,
		Content:  Content,
		TTL:      1,
		Priority: 10,
		Proxied:  Proxied,
	}

	// 认证字段
	header := getAuthHeader(config)

	body, err := json.Marshal(jsonData)

	request := &http.Request{
		Method: "POST",
		URL:    url,
		Header: header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(body)),
	}

	if err != nil {
		return result, err
	}

	response, err := myhttp.Client.Do(request)

	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	bodyResponse := createDNSRecordResponse{}

	err = json.NewDecoder(response.Body).Decode(&bodyResponse)

	if err != nil {
		return result, err
	}

	result = bodyResponse.Result

	return result, nil
}
