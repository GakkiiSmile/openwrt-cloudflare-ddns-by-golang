package api

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/services/myhttp"
	"net/http"
	"net/url"
)

type listDNSRecordResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Result []model.DNSRecord `json:"result"`
}

// GetDNSRecord 获取dns记录
func GetDNSRecord(config *model.Config, zoneID string, Type string) (list []model.DNSRecord, err error) {

	httpClient := myhttp.Client

	url, err := url.Parse(BASEURL + "/zones/" + zoneID + "/dns_records")

	if err != nil {
		return nil, err
	}

	url.Query().Add("type", Type)
	url.Query().Add("name", config.DomainName)
	url.Query().Add("order", "type")
	url.Query().Add("direction", "desc")
	url.Query().Add("match", "all")
	// url.Query().Add("page", "")
	// url.Query().Add("per_page", "")
	// url.Query().Add("content", "")

	// 认证字段
	header := getAuthHeader(config)

	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: header,
	}

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	bodyResponse := listDNSRecordResponse{}

	err = json.Unmarshal(body, &bodyResponse)

	if err != nil {
		return nil, err
	}

	return bodyResponse.Result, nil
}
