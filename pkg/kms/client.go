package kms

import (
	"aiot-service-for-mfp/pkg/kms/algorithms"
	"aiot-service-for-mfp/pkg/kms/config"
	"aiot-service-for-mfp/pkg/kms/exception"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Client struct {
	config *config.LongforConfig
}

func NewClientWithAccessKey(config *config.LongforConfig) *Client {
	return &Client{
		config: config,
	}
}

func (client *Client) DoAction(request LongforRequest, response interface{}) (err error) {
	timeStamp := time.Now().UTC().Format(time.RFC3339)
	request.addQueryParam("timeStamp", timeStamp)
	request.addQueryParam("accessKeyId", client.config.AccessKeyId)

	var params []string
	for k, v := range request.getQueryParams() {
		params = append(params, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	sort.Strings(params)

	var source string
	for _, v := range params {
		if len(source) > 0 {
			source += "&"
		}
		source = source + v
	}
	sham := url.QueryEscape(algorithms.ShaHmac1(source, client.config.AccessKeySecret))

	fullUrl := fmt.Sprintf("%s%s?%s&signature=%s", client.config.BasePath, request.getRoute(), source, sham)

	jsonData, err := json.Marshal(request.getBodyParam())
	if err != nil {
		return err
	}
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	if len(client.config.GaiaApiKey) > 0 {
		req.Header.Set("X-Gaia-Api-Key", client.config.GaiaApiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	err = client.readError(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, response)

}
func (client *Client) readError(data []byte) error {
	if len(data) == 0 {
		return errors.New("SDK.InvalidServerResponse")
	}
	var e = exception.HttpError{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}
	if e.Code > 0 && len(e.Message) > 0 {
		return errors.New(e.Message)
	}
	return nil
}
