package redis

import (
	"aiot-service-for-mfp/config"
	"aiot-service-for-mfp/pkg/logger"
	"fmt"
	//"github.com/garyburd/redigo/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	redigolib "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

func Connect() (redigolib.Conn, error) {

	setdb := redigolib.DialDatabase(int(config.RemoteConfs.RedisConf.DBName))
	setPasswd := redigolib.DialPassword(config.RemoteConfs.RedisConf.Password)

	c, err := redigolib.Dial("tcp", config.RemoteConfs.RedisConf.Host, setdb, setPasswd)
	if err != nil {
		logger.Logger.Error("Connect Redis Error1:", zap.Error(err))
		return nil, err
	}

	return c, nil

}

var RedisConn redigolib.Conn
var RedisPool *redigolib.Pool
var RedisSync *redsync.Redsync

func RedisPoolInit() *redigolib.Pool {
	return &redigolib.Pool{
		MaxIdle:     5,   //最大空闲数
		MaxActive:   500, //最大连接数，0不设上
		Wait:        true,
		IdleTimeout: time.Duration(1) * time.Second, //空闲等待时间
		Dial: func() (redigolib.Conn, error) {
			c, err := redigolib.Dial(
				"tcp",
				config.RemoteConfs.RedisConf.Host,
				redigolib.DialDatabase(int(config.RemoteConfs.RedisConf.DBName)),
				redigolib.DialPassword(config.RemoteConfs.RedisConf.Password), //redis IP地址
				redigolib.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redigolib.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redigolib.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
			)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigolib.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func RedisInit() {
	RedisPool = RedisPoolInit()
	RedisSync = redsync.New(redigo.NewPool(RedisPool))
	//RedisConn = RedisPool.Get()
	//mutex := RedisSync.NewMutex("test-redsync")
	//
	//if err := mutex.Lock(); err != nil {
	//	panic(err)
	//}
	//
	//if _, err := mutex.Unlock(); err != nil {
	//	panic(err)
	//}
}

func RedisClose() {
	RedisPool.Close()
	//RedisConn.Close()
}
