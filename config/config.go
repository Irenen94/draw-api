//@program: inference_agent
//@package: utils
//@file: config.go
//@description:配置文件加载
//@author: zhang.yuanhao
//@create: 2021-07-27 09:47

package config

import (
	"bytes"
	"aigc-image/internal/model"
	"aigc-image/pkg/logger"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

var (
	v     *viper.Viper
	Confs model.Conf
)

func init() {

	//获取默认设置
	box := packr.New("", "./conf")
	configType := "yaml"
	defaultConfig, _ := box.Find("config.yaml")

	v = viper.New()

	v.SetConfigType(configType)
	err := v.ReadConfig(bytes.NewReader(defaultConfig))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Load Config Err:%+v", err))
		panic(fmt.Sprintf("Load Config Err:%+v", err))
	}

	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, val := range configs {
		v.SetDefault(k, val)
	}

	//根据环境变量读取配置文件
	run_env := os.Getenv("ENV")
	if run_env == "" {
		run_env = "local"
	}

	fmt.Println("RUN ENV:", run_env)

	conf_name := "config"

	conf_name_builder := strings.Builder{}
	conf_name_builder.WriteString(conf_name)
	conf_name_builder.WriteString("-")
	conf_name_builder.WriteString(run_env)
	conf_name_builder.WriteString(".yaml")

	fmt.Println("Load Config:", conf_name_builder.String())

	envConfig, err := box.Find(conf_name_builder.String())

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Find EnvConfig Err:%+v", err))
	}

	err = v.ReadConfig(bytes.NewReader(envConfig))

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Load Env Config Err:%+v", err))
	}

	//序列化配置文件
	v.Unmarshal(&Confs)

	//根据运行环境设置ginmod
	Confs.BaseConf.RunMode = run_env
	if run_env == "local" {
		Confs.BaseConf.GinMode = "debug"
	} else {
		Confs.BaseConf.GinMode = "release"
	}
}
