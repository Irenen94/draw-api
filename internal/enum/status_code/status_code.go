package status_code

type StatusCode int

// 权限相关状态 20[0-9]
const (
	OK StatusCode = 200 + iota
	ERROR
	LOGIN_ERROR
	ACCESS_ERROR
	SIGN_EXPIRATION
	REPEAT_ERROR
	SIGN_ERROR
	TOKEN_ERROR
	TOKEN_FAIL
	THIRD_API_FAIL
	REQUEST_TOO_MANY
)

// 接口相关状态
const (
	REQUEST_ERROR          StatusCode = 211
	INCOMPLETE_REQUEST     StatusCode = 212
	DATA_PROCESS_ERROR     StatusCode = 213
	DATABASE_OPERATE_ERROR StatusCode = 214
)

//文件相关状态
const (
	FILE_HAVE_NO_DATA       StatusCode = 216
	FILE_TYPE_NOT_SUPPORT   StatusCode = 217
	IMG_DOWNLOAD_ERROR      StatusCode = 218
	IMG_ERROR               StatusCode = 219
	IMG_SIZE_SMALL_THAN_TWO StatusCode = 224
)

//人员信息相关状态
const (
	USER_NOT_EXISTS       StatusCode = 220
	USER_NOT_MATCH        StatusCode = 221
	USER_EXISTS           StatusCode = 222
	ACCOUNTTYPE_NOT_EXIST StatusCode = 223
)

// 认证ID
const (
	CERTIFY_NOT_EXIST         StatusCode = 230
	ALIYUN_CERTIFY_NOT_ACCESS StatusCode = 231
	VERIFY_TYPE_UNVALID       StatusCode = 232
	SELFALGO_UNVALID          StatusCode = 235
	LFFACE_TOKEN_NOT_EXIST    StatusCode = 237
	NO_SCENE_THRESHOLD        StatusCode = 238
)

// 算法结果
const (
	MULTI_COMPARE_FAIL StatusCode = 1001
	IDENTIFY_FAIL      StatusCode = 1002
	FACE_COMPARE_FAIL  StatusCode = 1003
)

// redis
const (
	REDIS_OPERATE_ERROR StatusCode = 240
)

const (
	NO_IMAGE_DATA      StatusCode = 272
	INSERT_PERSON_FAIL StatusCode = 273
)

func (c StatusCode) String() string {
	switch c {
	case OK:
		return "成功"
	case ERROR:
		return "失败"
	case LOGIN_ERROR:
		return "用户名或密码错误"
	case ACCESS_ERROR:
		return "权限不足，请鉴权"
	case SIGN_EXPIRATION:
		return "签名已过期"
	case REPEAT_ERROR:
		return "数据重复"
	case SIGN_ERROR:
		return "签名验证未通过，鉴权失败"
	case TOKEN_ERROR:
		return "token为空"
	case TOKEN_FAIL:
		return "token验证失败，请重新获取"
	case THIRD_API_FAIL:
		return "算法接口调用失败"
	case REQUEST_ERROR:
		return "请求数据异常"
	case INCOMPLETE_REQUEST:
		return "请求数据不全，或者字段不合规"
	case REQUEST_TOO_MANY:
		return "当前请求过多，请稍后再试"
	case DATA_PROCESS_ERROR:
		return "数据处理错误"
	case DATABASE_OPERATE_ERROR:
		return "数据库操作失败"
	case FILE_HAVE_NO_DATA:
		return "文件无有效数据"
	case FILE_TYPE_NOT_SUPPORT:
		return "文件类型暂不支持"
	case IMG_DOWNLOAD_ERROR:
		return "图片下载失败"
	case IMG_ERROR:
		return "图片处理错误"
	case IMG_SIZE_SMALL_THAN_TWO:
		return "img size small than two"
	case USER_NOT_EXISTS:
		return "用户信息不存在"
	case USER_NOT_MATCH:
		return "用户信息不匹配"
	case USER_EXISTS:
		return "用户信息已存在"
	case ACCOUNTTYPE_NOT_EXIST:
		return "账户类型不存在"
	case NO_IMAGE_DATA:
		return "服务未查得对应底图信息"
	case CERTIFY_NOT_EXIST:
		return "认证ID不存在"
	case REDIS_OPERATE_ERROR:
		return "redis操作失败"
	case ALIYUN_CERTIFY_NOT_ACCESS:
		return "aliyun认证ID没有此接口权限"
	case MULTI_COMPARE_FAIL:
		return "多张人脸比对失败"
	case IDENTIFY_FAIL:
		return "aliyun实人认证失败"
	case VERIFY_TYPE_UNVALID:
		return "verifyType类型有误"
	case INSERT_PERSON_FAIL:
		return "数据写入人脸库失败"
	case SELFALGO_UNVALID:
		return "selfAlgo类型有误"
	case FACE_COMPARE_FAIL:
		return "上传人脸与底库比对失败"
	case LFFACE_TOKEN_NOT_EXIST:
		return "lffaceToken无效"
	default:
		return ""
	}
}
