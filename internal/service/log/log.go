package log

import (
	"go.uber.org/zap"
	"aigc-image/config"
	"aigc-image/internal/model"
	"aigc-image/pkg/logger"
	"aigc-image/pkg/mysql"
)

type algoLog struct{}

var AlgoLogService algoLog

func (q *algoLog) CreateAlgoLog(log model.AlgoLog) (err error) {

	db := mysql.Mysql.DBs[config.Confs.PolarConf.DBName]

	if db == nil {
		err = mysql.Mysql.InitMysqlDB(config.Confs.PolarConf)
		db = mysql.Mysql.DBs[config.Confs.PolarConf.DBName]

		if err != nil || db == nil {
			logger.Logger.Error("InitMysqlDB Err:", zap.Error(err))
			return err
		}
	}

	return db.Table("algo_log").Create(&log).Error
}
