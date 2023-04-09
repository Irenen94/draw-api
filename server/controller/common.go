package controller

import (
	"aigc-image/internal/enum/status_code"
	"aigc-image/internal/model"
	"aigc-image/internal/service/algo"
	"aigc-image/internal/service/log"
	"aigc-image/internal/util"
	"aigc-image/pkg/logger"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	_ "image/jpeg"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var Limiter = rate.NewLimiter(1,1)

func CreateImage(c *gin.Context) {
	if !Limiter.Allow(){
		c.JSON(http.StatusOK, gin.H{
			"requestId": "",
			"trace":     "",
			"msg":  status_code.REQUEST_TOO_MANY.String(),
			"code": status_code.REQUEST_TOO_MANY,
			"data": "-1",
		})
		return
	}
	requestId := c.Request.Header.Get("X-Request-ID")
	st := time.Now()
	var param model.RequestAlgo
	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"requestId": requestId,
			"trace":     err,
			"msg":  status_code.INCOMPLETE_REQUEST.String(),
			"code": status_code.INCOMPLETE_REQUEST,
			"data": "-1",
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(param)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusOK, gin.H{
				"requestId": requestId,
				"trace":     e.Error(),
				"message":   status_code.INCOMPLETE_REQUEST.String(),
				"code":      status_code.INCOMPLETE_REQUEST,
				"data":      "-1",
			})
			break
		}
		return
	}
	if param.SmartMode {
		param.Prompt = param.Prompt + ", pastel colors, highly detailed, digital painting, concept art"
		//	param.Prompt = param.Prompt + "," + "Concept art" + "," + "Thomas Moran"
		param.DdimSteps = 90
		//param.Seed = 87
	}
	param.Prompt = strings.Replace(strings.Replace(param.Prompt, "【", "[", -1), "】", "]", -1)
	var negative string
	param.SizeType = 2

	r := regexp.MustCompile(`\[([^]]+)\]`)
	matchs := r.FindAllStringSubmatch(param.Prompt, -1)
	for _, s := range matchs {
		fmt.Println(s[1])
		negative += s[1] + ","
	}
	//fmt.Println(matchs[0][1])]
	if len(param.Negative) != 0 {
		param.Negative = param.Negative + "," + negative
	} else {
		param.Negative = negative
	}

	sampleRegexp := regexp.MustCompile(`\<.*?>|\[.*?]`)
	param.Prompt = sampleRegexp.ReplaceAllString(param.Prompt, "")

	tokenCode, result := algo.AlgoService.CreateImage(requestId, &param, 100000)
	if tokenCode != 0 {
		logger.Logger.Error("StateSubmitErrorLog", zap.String("RequestId", requestId), zap.String("Method", "TokenValidateFromDb"), zap.Error(errors.New("token校验失败")))
		c.JSON(http.StatusOK, gin.H{
			"requestId": requestId,
			"trace":     status_code.THIRD_API_FAIL.String(),
			"message":   result.Msg,
			"code":      status_code.THIRD_API_FAIL,
			"data":      result.Data,
		})
		//return
	}

	totalStartTime := st.UnixNano() / 1e6
	totalEndTime := st.UnixNano() / 1e6
	stateLog := model.AlgoLog{
		RequestId: requestId,
		Prompt:    param.Prompt,
		TaskId:    param.TaskId,
		DdimSteps: param.DdimSteps,
		Negative:  param.Negative,
		SizeType:  param.SizeType,
		H:            param.H,
		W:            param.W,
		NSamples:     param.NSamples,
		ImageType:    "simple",
		SmartMode:    "false",
		Seed:         param.Seed,
		InitImage:    "",
		CreateImages: "",
		SnapTime:     st,
		Code:         result.Code,
		Message:      result.Msg,
		Details:      "",
		TotalTime:    0,
	}
	if param.SmartMode {
		stateLog.SmartMode = "true"
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		if len(param.InitImg) != 0 {
			stateLog.InitImage, err = algo.AlgoService.FileTOOss(string(param.InitImg), ".jpg")
			if err != nil {
				logger.Logger.Info("上传oss失败！ ", zap.Error(err))
			}
		}
	}()

	ossImages := make([]string, len(result.Data))
	if len(result.Data) > 0 {
		for i := range result.Data {
			wg.Add(1)
			go func(i int) {
				defer func() {
					wg.Done()
				}()
				var ossUrl string
				var err2 error
				task := func() error {
					ossUrl, err2 = algo.AlgoService.FileTOOss(string(result.Data[i]), ".jpg")
					if err2 != nil {
						logger.Logger.Error("LivenessSubmitImagesErrorLog", zap.String("RequestID:", requestId), zap.String("Method", "DealWithImagesWithoutAES"), zap.Error(err))
					}
					return err2
				}
				//设定3秒最大重试次数，最大3次，1秒重试一次
				err2 = util.RetryTimes("更新任务", 3, 5*time.Millisecond, task)
				if err2 != nil {
					logger.Logger.Error("CompareImageWithBase64", zap.Error(err2))
				}
				ossImages[i] = ossUrl
			}(i)
		}
	}

	wg.Wait()
	stateLog.CreateImages = strings.Join(ossImages, ";")

	// 日志入库
	defer func() {
		err = log.AlgoLogService.CreateAlgoLog(stateLog)
		if err != nil {
			logger.Logger.Error("StateSubmitErrorLog", zap.String("RequestID:", requestId), zap.String("Method", "CreateStateSubmitLog"), zap.Error(err))
		}
	}()

	defer func() {
		totalEndTime = time.Now().UnixNano() / 1e6
		totalTime := totalEndTime - totalStartTime
		stateLog.TotalTime = int(totalTime)
	}()

	stateLog.Code = int(status_code.OK)
	stateLog.Message = status_code.OK.String()
	c.JSON(http.StatusOK, gin.H{
		"requestId": requestId,
		"trace":     "",
		"message":   "成功",
		"code":      status_code.OK,
		"data":      result.Data,
	})
	return
}

func CreateProImage(c *gin.Context) {
	if !Limiter.Allow(){
		c.JSON(http.StatusOK, gin.H{
			"requestId": "",
			"trace":     "",
			"msg":  status_code.REQUEST_TOO_MANY.String(),
			"code": status_code.REQUEST_TOO_MANY,
			"data": "-1",
		})
		return
	}
	requestId := c.Request.Header.Get("X-Request-ID")
	st := time.Now()
	var param model.RequestProAlgo
	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  status_code.INCOMPLETE_REQUEST.String(),
			"code": status_code.INCOMPLETE_REQUEST,
			"data": "-1",
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(param)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusOK, gin.H{
				"requestId": requestId,
				"trace":     e.Error(),
				"message":   status_code.INCOMPLETE_REQUEST.String(),
				"code":      status_code.INCOMPLETE_REQUEST,
				"data":      "-1",
			})
			break
		}
		return
	}

	var paramList []model.RequestAlgo
	if len(param.Artists) == 0 {
		param.Artists = "origin"
	}

	if len(param.Styles) == 0 {
		param.Styles = "origin"
	}
	totalStartTime := st.UnixNano() / 1e6
	totalEndTime := st.UnixNano() / 1e6
	var artists = strings.Split(param.Artists, ",")
	var styles = strings.Split(param.Styles, ",")

	if param.Strength < 0.01 {
		param.Strength = 1
	}
	if param.DdimSteps < 20 {
		param.DdimSteps = 20
	}
	if param.DdimSteps > 99 {
		param.DdimSteps = 99
	}
	if param.Seed == 0 {
		param.Seed = algo.GenerateRandInt(0, 100)
	}
	if param.NSamples < 1 {
		param.NSamples = 1
	}
	if param.NSamples > 8 {
		param.NSamples = 8
	}

	param.Prompt = strings.Replace(strings.Replace(param.Prompt, "【", "[", -1), "】", "]", -1)
	var negative string

	r := regexp.MustCompile(`\[([^]]+)\]`)
	matchs := r.FindAllStringSubmatch(param.Prompt, -1)
	for _, s := range matchs {
		fmt.Println(s[1])
		negative += s[1] + ","
	}
	//fmt.Println(matchs[0][1])]
	if len(param.Negative) != 0 {
		param.Negative = param.Negative + "," + negative
	} else {
		param.Negative = negative
	}

	sampleRegexp := regexp.MustCompile(`\<.*?>|\[.*?]`)
	param.Prompt = sampleRegexp.ReplaceAllString(param.Prompt, "")

	for i := range artists {
		for j := range styles {
			var tmpParam model.RequestAlgo
			tmpParam.Artist = artists[i]
			tmpParam.Style = styles[j]
			tmpParam.Seed = param.Seed
			tmpParam.Negative = param.Negative
			tmpParam.Prompt = param.Prompt + "," + fmt.Sprintf("%v", styles[j]) + "," + fmt.Sprintf("%v", artists[i])
			tmpParam.Strength = param.Strength
			tmpParam.TaskId = param.TaskId
			tmpParam.InitImg = param.InitImg
			tmpParam.DdimSteps = param.DdimSteps
			tmpParam.SizeType = param.SizeType
			//tmpParam.H = param.H
			//tmpParam.W = param.W
			tmpParam.TaskId = param.TaskId
			tmpParam.NSamples = 1
			paramList = append(paramList, tmpParam)
		}
	}

	stateLog := model.AlgoLog{
		RequestId: requestId,
		Prompt:    param.Prompt,
		TaskId:    param.TaskId,
		DdimSteps: param.DdimSteps,
		Negative:  param.Negative,
		Styles:    param.Styles,
		Artists:   param.Artists,
		SizeType:  param.SizeType,
		//H:            param.H,
		//W:            param.W,
		NSamples:     param.NSamples,
		ImageType:    "pro",
		Seed:         param.Seed,
		InitImage:    "",
		CreateImages: "",
		SnapTime:     st,
		Code:         200,
		Message:      "",
		Details:      "",
		TotalTime:    0,
	}

	// 日志入库
	defer func() {
		err = log.AlgoLogService.CreateAlgoLog(stateLog)
		if err != nil {
			logger.Logger.Error("StateSubmitErrorLog", zap.String("RequestID:", requestId), zap.String("Method", "CreateStateSubmitLog"), zap.Error(err))
		}
	}()

	defer func() {
		totalEndTime = time.Now().UnixNano() / 1e6
		totalTime := totalEndTime - totalStartTime
		stateLog.TotalTime = int(totalTime)
	}()

	//sets := Product([]interface{}{artists}, []interface{}{styles})

	wg := sync.WaitGroup{}

	results := make([]model.ResponseProAlgo, len(paramList))
	var result model.ResponseProAlgo
	var tokenCode int32
	for i := range paramList {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			tokenCode, result = algo.AlgoService.CreateImage(requestId, &paramList[i], 100000)
			if tokenCode != 0 {
				logger.Logger.Error("LivenessSubmitImagesErrorLog", zap.String("RequestID:", requestId), zap.String("Method", "DealWithImagesWithoutAES"), zap.Error(err))
			}
			results[i] = result
		}(i)
	}
	wg.Wait()

	stateLog.Code = result.Code
	stateLog.Message = result.Msg

	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go func() {

		defer func() {
			wg1.Done()
		}()
		if len(param.InitImg) != 0 {
			stateLog.InitImage, err = algo.AlgoService.FileTOOss(string(param.InitImg), ".jpg")
			if err != nil {
				logger.Logger.Info("上传oss失败！ ", zap.Error(err))
			}
		}
	}()

	ossImages := make([]string, len(results))
	for i := range results {
		wg1.Add(1)
		go func(i int) {
			defer func() {
				wg1.Done()
			}()
			var ossUrl string
			//var err error
			if len(results[i].Data) != 0 {
				ossUrl, err = algo.AlgoService.FileTOOss(results[i].Data[0], ".jpg")
				if err != nil {
					logger.Logger.Info("上传oss失败！ ", zap.Error(err))
				}
			}
			ossImages[i] = ossUrl
		}(i)
	}

	wg1.Wait()

	var resultList model.ResponseProAlgoList
	resultList.RequestId = requestId
	resultList.Code = 200
	resultList.Msg = "success"
	for i := range results {
		if results[i].Code != 200 {
			resultList.Code = results[i].Code
			resultList.Msg = results[i].Msg
		}
		var tmpData model.ResponseImage
		tmpData.Artist = results[i].Artist
		tmpData.Style = results[i].Style
		if len(results[i].Data) > 0 {
			tmpData.Image = results[i].Data[0]
		}
		resultList.Data = append(resultList.Data, tmpData)
	}

	stateLog.CreateImages = strings.Join(ossImages, ";")
	stateLog.Code = resultList.Code
	stateLog.Message = resultList.Msg
	//stateLog.Details = err.Error()
	if resultList.Code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"requestId": requestId,
			"trace":     status_code.THIRD_API_FAIL.String(),
			"message":   resultList.Msg,
			"code":      status_code.THIRD_API_FAIL,
			"data":      resultList.Data,
		})
		return
	}

	stateLog.Code = int(status_code.OK)
	stateLog.Message = status_code.OK.String()
	c.JSON(http.StatusOK, gin.H{
		"requestId": requestId,
		"trace":     "",
		"message":   "成功",
		"code":      status_code.OK,
		"data":      resultList.Data,
	})
	return
}

//笛卡尔积算法
func Product(sets ...[]interface{}) [][]interface{} {
	lens := func(i int) int { return len(sets[i]) }
	product := [][]interface{}{}
	for ix := make([]int, len(sets)); ix[0] < lens(0); nextIndex(ix, lens) {
		var r []interface{}
		for j, k := range ix {
			r = append(r, sets[j][k])
		}
		product = append(product, r)
	}
	return product
}

func nextIndex(ix []int, lens func(i int) int) {
	for j := len(ix) - 1; j >= 0; j-- {
		ix[j]++
		if j == 0 || ix[j] < lens(j) {
			return
		}
		ix[j] = 0
	}
}

func ListRecords(c *gin.Context) {
	requestId := c.Request.Header.Get("X-Request-ID")
	//var param model.RequestAlgo
	//err := c.ShouldBindBodyWith(&param, binding.JSON)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg":  status_code.INCOMPLETE_REQUEST.String(),
	//		"code": status_code.INCOMPLETE_REQUEST,
	//		"data": "-1",
	//	})
	//	return
	//}
	//
	//validate := validator.New()
	//err = validate.Struct(param)
	//if err != nil {
	//	for _, e := range err.(validator.ValidationErrors) {
	//		c.JSON(http.StatusOK, gin.H{
	//			"requestId": requestId,
	//			"trace":     e.Error(),
	//			"message":   status_code.INCOMPLETE_REQUEST.String(),
	//			"code":      status_code.INCOMPLETE_REQUEST,
	//			"data":      "-1",
	//		})
	//		break
	//	}
	//	return
	//}

	result, err := algo.AlgoService.ListRecords(requestId)
	if err != nil {
		logger.Logger.Error("StateSubmitErrorLog", zap.String("RequestId", requestId), zap.String("Method", "TokenValidateFromDb"), zap.Error(errors.New("token校验失败")))
		c.JSON(http.StatusOK, gin.H{
			"requestId": requestId,
			"trace":     err,
			"message":   "失败",
			"code":      500,
			"data":      "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"code":    status_code.OK,
		"data":    result,
	})
	return
}
