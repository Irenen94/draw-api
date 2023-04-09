//@program: faceVerify
//@package: model
//@file: config.go
//@description:
//@author: shiyajie
//@create: 2022-04-11 10:24

package model

type Conf struct {
	BaseConf  BaseConfig  `yaml:"baseConf"`
	PolarConf PolarConfig `yaml:"polarConf"`
	OSSConf   OSSConf     `yaml:"ossConf"`
	RpcConf   RpcConfig   `yaml:"rpcConf"`
}

type BaseConfig struct {
	RunMode         string `yaml:"runMode"`      //运行环境
	GinMode         string `yaml:"ginMode"`      //gin 模式
	HTTPPort        int    `yaml:"httpPort"`     //http端口
	ReadTimeout     int    `yaml:"readTimeout"`  //读取超时时间
	WriteTimeout    int    `yaml:"writeTimeout"` //写入超时时间
	SceneId         int64  `yaml:"scene_id"`
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	AlgoPath        string `yaml:"algoPath"`
}

type RpcConfig struct {
	IDMAddr string `yaml:"idmAddr"`
}

type PolarConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DBName   string `yaml:"dbName"`
}

type OSSConf struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
}
