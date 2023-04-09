package algo

import (
	"aigc-image/config"
	"aigc-image/internal/model"
	"aigc-image/internal/util"
	"aigc-image/pkg/logger"
	"aigc-image/pkg/mysql"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"strings"
	"time"
)

type algoService struct {
}

var AlgoService algoService

func (self *algoService) CreateImage(requestId string, param *model.RequestAlgo, timeout int) (code int32, result model.ResponseProAlgo) {

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
		param.Seed = GenerateRandInt(0, 100)
	}
	if param.NSamples < 1 {
		param.NSamples = 1
	}
	if param.NSamples > 8 {
		param.NSamples = 8
	}
	result.Artist = param.Artist
	result.Style = param.Style
	result.Data = []string{}
	result.Msg = "error"
	result.Code = 500
	req, _ := json.Marshal(param)
	headerMap := map[string]string{}
	ok, code, res, _ := util.CallServiceWithBytes(config.Confs.BaseConf.AlgoPath, req, "POST", headerMap, timeout)

	if ok && code == 200 {
		if e := json.Unmarshal(res, &result); e == nil {
			if result.Code != 200 {
				return 1001, result
			}
			return 0, result
		} else {
			logger.Logger.Error("gpu status json err", zap.String("err", e.Error()))
			return 1000, result
		}
	}
	return code, result
}

func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func GenerateRandInt(min, max int) int {
	rand.Seed(time.Now().Unix()) //随机种子
	return rand.Intn(max-min) + min
}

func (self *algoService) FileTOOss(imageBase64 string, suffix string) (result string, err error) {
	videoMd5 := MD5(imageBase64)
	imageBase64Byte, _ := base64.StdEncoding.DecodeString(imageBase64)

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(config.Confs.OSSConf.Endpoint, config.Confs.OSSConf.AccessKeyID, config.Confs.OSSConf.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(config.Confs.OSSConf.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 指定Object存储类型为低频访问。
	storageType := oss.ObjectStorageClass(oss.StorageIA)

	// 指定Object访问权限为私有。
	//objectAcl := oss.ObjectACL(oss.ACLPrivate)

	// 将字符串"Hello OSS"上传至exampledir目录下的exampleobject.txt文件。
	fileName := "image-dragon/" + videoMd5 + suffix
	err = bucket.PutObject(fileName, strings.NewReader(string(imageBase64Byte)), storageType)
	result, err = bucket.SignURL(fileName, oss.HTTPGet, 3600*24*365*100)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println(result)
	return
}

func (self *algoService) ListRecords(requestId string) (result []model.AlgoLog, err error) {

	db := mysql.Mysql.DBs[config.Confs.PolarConf.DBName]
	if db == nil {
		err = mysql.Mysql.InitMysqlDB(config.Confs.PolarConf)
		db = mysql.Mysql.DBs[config.Confs.PolarConf.DBName]
		if err != nil || db == nil {
			logger.Logger.Error("TokenValidateFromDbErrorLog", zap.String("RequestID:", requestId), zap.String("Method", "InitMysqlDB"), zap.Error(err))
			return result, err
		}
	}
	data := []model.AlgoLog{}
	// 防注入
	err = db.Table("algo_log").Scan(&data).Error
	if err != nil {
		logger.Logger.Error("TokenValidateFromDbErrorLog", zap.String("RequestID:", requestId), zap.String("Method", "DBScan"), zap.Error(err))
		return result, err
	}
	//if len(result) == 0 {
	//	logger.Logger.Error("TokenValidateFromDbErrorLog", zap.String("RequestID:", requestId), zap.Error(errors.New(status_code.LFFACE_TOKEN_NOT_EXIST.String())))
	//	return result, err
	//}
	return data, err
}
