package exception

type HttpError struct {
	RequestId string `json:"requestId"` // 本次请求的ID
	Code      int    `json:"code"`
	Message   string `json:"message"`
}
