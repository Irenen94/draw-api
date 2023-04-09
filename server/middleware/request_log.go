//@program: unified_face_datastore
//@package: middleware
//@file: request_log.go
//@description:
//@author: zhang.yuanhao
//@create: 2022-04-19 17:34

package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
	"aigc-image/pkg/logger"
)

type GatewayLog struct {
	RequestID   string `json:"requestid"` //请求id
	ReqTime     int64  `json:"req_time"`  //请求时间
	Url         string `json:"url"`
	Host        string `json:"host"`
	Path        string `json:"path"`          //请求地址
	Method      string `json:"method"`        //请求方法
	ReqBodySize int64  `json:"req_body_size"` //请求体大小
	ResBodySize int    `json:"res_body_size"` //返回体大小
	ResCode     int    `json:"res_code"`      //返回状态码
	ResTime     int64  `json:"res_time"`      //响应时常
	LogTime     int64  `json:"log_time"`      //记录时间
}

// 自定义结构并组合实现 Web-Gin.ResponseWrite 接口
type MyResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Body 获取复写 body
func (rw *MyResponseWriter) Body() *bytes.Buffer {
	return rw.body
}

// Write 重写 Write 方法，将 response 重复写到 rw.body
func (rw *MyResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// WriteString 重写 WriteString 方法，将 response 重复写到 rw.body
func (rw *MyResponseWriter) WriteString(s string) (int, error) {
	rw.body.WriteString(s)
	return rw.ResponseWriter.WriteString(s)
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		st := time.Now()

		//ctx Writer 重写
		rw := &MyResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = rw
		c.Next()

		go func(c *gin.Context, s time.Time) {
			url := c.Request.URL.String()
			path := c.Request.URL.Path
			host := c.Request.Host
			method := c.Request.Method
			requestId := c.Request.Header.Get("X-Request-ID")

			log := GatewayLog{
				RequestID:   requestId,
				ReqTime:     st.UnixNano() / 1e6,
				Url:         url,
				Path:        path,
				Host:        host,
				Method:      method,
				ReqBodySize: c.Request.ContentLength,
				ResBodySize: c.Writer.Size(),
				ResCode:     c.Writer.Status(),
				ResTime:     time.Now().Sub(st).Milliseconds(),
				LogTime:     time.Now().UnixNano() / 1e6,
			}

			if !strings.Contains(rw.Body().String(), "\"code\":200") && log.ResBodySize < 10000 {
				logger.Logger.Info("RequestErrorLog", zap.String("RequestID:", requestId), zap.Any("Log:", log), zap.String("Response:", rw.Body().String()))
			}

			return

		}(c, st)

	}
}
