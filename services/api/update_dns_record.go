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

// UpdateDNSRecordResponse 响应结果
type updateDNSRecordResponse struct {
	Success  bool            `json:"success"`
	Errors   []interface{}   `json:"errors"`
	Messages []interface{}   `json:"messages"`
	Result   model.DNSRecord `json:"result"`
}

type postUpdateDNSJSONData struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// UpdateDNSRecord 获取dns记录
func UpdateDNSRecord(config *model.Config, id string, Type string, Name string, Content string, Proxied bool) (result model.DNSRecord, err error) {

	httpClient := myhttp.Client

	url, err := url.Parse(BASEURL + "/zones/" + config.ZoneID + "/dns_records/" + id)
	result = model.DNSRecord{}

	if err != nil {
		return result, err
	}

	postData := &postUpdateDNSJSONData{
		Type:    Type,
		Name:    Name,
		Content: Content,
		TTL:     1,
		Proxied: Proxied,
	}

	data, _ := json.Marshal(postData)

	// 认证字段
	header := getAuthHeader(config)

	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(data)),
	}

	response, err := httpClient.Do(request)

	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	bodyResponse := updateDNSRecordResponse{}

	err = json.Unmarshal(body, &bodyResponse)

	if err != nil {
		return result, err
	}

	result = bodyResponse.Result

	return result, nil
}
