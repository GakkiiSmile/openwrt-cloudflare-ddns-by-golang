package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"
	"main/services/myhttp"
	"net/http"
	"net/url"
	"time"
)

// Zone 结构体
type Zone struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	Status              string      `json:"status"`
	Paused              bool        `json:"paused"`
	Type                string      `json:"type"`
	DevelopmentMode     int         `json:"development_mode"`
	NameServers         []string    `json:"name_servers"`
	OriginalNameServers []string    `json:"original_name_servers"`
	OriginalRegistrar   interface{} `json:"original_registrar"`
	OriginalDnshost     interface{} `json:"original_dnshost"`
	ModifiedOn          time.Time   `json:"modified_on"`
	CreatedOn           time.Time   `json:"created_on"`
	ActivatedOn         time.Time   `json:"activated_on"`
}

// ZoneResponse 相应结构体
type zoneResponse struct {
	Result     []Zone `json:"result"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		TotalPages int `json:"total_pages"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
	} `json:"result_info"`
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// GetZoneRecord 获取zone
func GetZoneRecord(config *model.Config) (list []Zone, err error) {

	httpClient := myhttp.Client

	url, err := url.Parse(BASEURL + "/zones")

	if err != nil {
		return nil, err
	}

	url.Query().Add("name", config.DomainName)
	url.Query().Add("order", "name")
	url.Query().Add("direction", "desc")
	url.Query().Add("match", "all")

	// 认证字段
	header := getAuthHeader(config)

	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: header,
	}

	response, err := httpClient.Do(request)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	bodyResponse := zoneResponse{}

	err = json.Unmarshal(body, &bodyResponse)

	if err != nil {
		return nil, err
	}

	return bodyResponse.Result, nil
}
