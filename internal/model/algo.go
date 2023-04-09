package model

//type RequestAlgo struct {
//	SceneId   int64    `json:"sceneId" form:"sceneId" binding:"required"`
//	CertifyId string   `json:"certifyId" form:"certifyId" binding:"required"`
//	Test      bool     `json:"test" form:"test" defaults:"false"`
//	CatchOff  bool     `json:"catchOff" form:"catchOff" default:"false"` //是否关闭缓存数据，默认不关闭
//	Images    []string `json:"images"`
//}

// 简单模式
type RequestAlgo struct {
	Prompt    string  `json:"prompt" binding:"required"`
	TaskId    string  `json:"task_id"`
	DdimSteps int     `json:"ddim_steps" default:"30"`
	Strength  float64 `json:"strength"`
	Negative  string  `json:"negative"`
	SizeType  int     `json:"size_type"`
	H         int     `json:"H" default:"512"`
	W         int     `json:"W" default:"512"`
	NSamples  int     `json:"n_samples" default:"1"`
	Seed      int     `json:"seed"`
	InitImg   string  `json:"init_img"`
	SmartMode bool    `json:"smart_mode"` //智能模式是否开启，true开启，false关闭，默认关闭
	Artist    string  `json:"artist"`
	Style     string  `json:"style"`
}

// 专业模式输入
type RequestProAlgo struct {
	Prompt    string  `json:"prompt" binding:"required"`
	TaskId    string  `json:"task_id"`
	DdimSteps int     `json:"ddim_steps"`
	Negative  string  `json:"negative"`
	Strength  float64 `json:"strength"`
	Artists   string  `json:"artists"`
	Styles    string  `json:"styles"`
	SizeType  int     `json:"size_type"`
	//H         int     `json:"H" default:"512"`
	//W         int     `json:"W" default:"512"`
	NSamples  int    `json:"n_samples"`
	Seed      int    `json:"seed"`
	InitImg   string `json:"init_img"`
	SmartMode bool   `json:"smart_mode"` //智能模式是否开启，true开启，false关闭，默认关闭
}

type ResponseAlgo struct {
	RequestId string   `json:"requestId"`
	Trace     string   `json:"trace"`
	Code      int      `json:"code"`
	Msg       string   `json:"msg"`
	Data      []string `json:"data"`
}

type ResponseProAlgo struct {
	RequestId string   `json:"requestId"`
	Trace     string   `json:"trace"`
	Artist    string   `json:"artist"`
	Style     string   `json:"style"`
	Code      int      `json:"code"`
	Msg       string   `json:"msg"`
	Data      []string `json:"data"`
}

type ResponseProAlgoList struct {
	RequestId string          `json:"requestId"`
	Trace     string          `json:"trace"`
	Code      int             `json:"code"`
	Msg       string          `json:"msg"`
	Data      []ResponseImage `json:"data"`
}

type ResponseImage struct {
	Artist string `json:"artist"`
	Style  string `json:"style"`
	Image  string `json:"image"`
}
