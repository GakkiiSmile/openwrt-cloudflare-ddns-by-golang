package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"
	"main/services/myhttp"
	"net/http"
	"net/url"
)

// DeleteResponse 响应结果
type deleteResponse struct {
	Result struct {
		ID string `json:"id"`
	} `json:"result"`
}

// DeleteDNSRecord 获取dns记录
func DeleteDNSRecord(config *model.Config, ID string) (id string, err error) {

	httpClient := myhttp.Client

	url, err := url.Parse(BASEURL + "/zones/" + config.ZoneID + "/dns_records/" + ID)

	if err != nil {
		return id, err
	}

	// 认证字段
	header := getAuthHeader(config)

	request := &http.Request{
		Method: "DELETE",
		URL:    url,
		Header: header,
	}

	response, err := httpClient.Do(request)

	if err != nil {
		fmt.Println(err)
		return id, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	bodyResponse := deleteResponse{}

	err = json.Unmarshal(body, &bodyResponse)

	if err != nil {
		return id, err
	}

	id = bodyResponse.Result.ID

	return id, nil
}
