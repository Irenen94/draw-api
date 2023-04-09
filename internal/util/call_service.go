package util

import (
	"aigc-image/pkg/logger"
	"bytes"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

func CallServiceWithBytes(url string, req []byte, method string, headers map[string]string, timeout int) (bool, int32, []byte, string) {
	//logger.Logger.Debug("Call service.", zap.String("url", url))
	client := &http.Client{
		Timeout: time.Millisecond * time.Duration(timeout), // Set ms timeout.
	}
	//如果失败， retry
	retryTimes := 3
	var response *http.Response = nil
	var err1 error
	var succeed = true
	// 这里的重试是无效的
	// 会发生如下bug {"err": "Post \"http://127.0.0.1:10000/cv/generic/person_body_detection\": http: ContentLength=298698 with Body length 0"}
	//按照这里的说法， 每次重试需要重建请求 https://blog.csdn.net/weixin_30430169/article/details/96396430?utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromMachineLearnPai2%7Edefault-1.baidujs&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromMachineLearnPai2%7Edefault-1.baidujs
	for i := 0; i < retryTimes; i++ {
		if !succeed {
			logger.Logger.Info("Retry do request.", zap.Int("seq", i))
		}
		reader := bytes.NewReader(req)
		request, err := http.NewRequest(method, url, reader)
		if err != nil {
			logger.Logger.Error("Err when build request.", zap.String("err", err.Error()))
			return false, 0, nil, "build request failed"
		}
		for k, v := range headers {
			request.Header.Set(k, v)
		}
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")

		response, err1 = client.Do(request)
		if err1 != nil {
			succeed = false
			logger.Logger.Error("Err when post http request.", zap.String("err", err1.Error()))
		} else {
			succeed = true
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
	} else {
		return false, 0, nil, "no response"
	}
	if !succeed {
		logger.Logger.Error("Failed call service.", zap.Int("retry_times", retryTimes))
		return false, 0, nil, "call service failed"
	}
	if response.StatusCode == 404 {
		return false, 404, nil, "Service not exist"
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Logger.Error("Read response body error", zap.String("err", err.Error()))
		return false, int32(response.StatusCode), nil, "Read response body error"
	}
	return true, int32(response.StatusCode), body, "ok"
}
