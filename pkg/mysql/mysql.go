package mysql

import (
	"errors"
	"fmt"
	"time"

	"aigc-image/internal/model"
	"aigc-image/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultMaxIdleConn     = 10
	defaultMaxOpenConn     = 50
	defaultConnMaxLifetime = 3 * time.Hour
)

type mysql_db struct {
	DBs map[string]*gorm.DB
}

var Mysql = new(mysql_db)

func (*mysql_db) InitMysqlDB(cfg ...model.PolarConfig) (err error) {
	dbs, err := NewMulti(WithConfigs(cfg...))
	if err != nil {
		return
	}
	Mysql.DBs = dbs
	return
}

type Option func(*option)

type option struct {
	dbConfigs       []model.PolarConfig
	gormConfig      gorm.Config
	maxIdleConn     int           // 空闲连接池中连接的最大数量
	maxOpenConn     int           // 打开数据库连接的最大数量
	connMaxLifetime time.Duration // 连接可复用的最大时间
}

func WithConfigs(cfg ...model.PolarConfig) Option {
	return func(o *option) {
		o.dbConfigs = cfg
	}
}

func WithGormConfig(cfg gorm.Config) Option {
	return func(o *option) {
		o.gormConfig = cfg
	}
}

func WithMaxIdleConn(maxIdleConn int) Option {
	return func(o *option) {
		o.maxIdleConn = maxIdleConn
	}
}

func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(o *option) {
		o.connMaxLifetime = connMaxLifetime
	}
}

func WithMaxOpenConn(maxOpenConn int) Option {
	return func(o *option) {
		o.maxOpenConn = maxOpenConn
	}
}

// New 初始化一个数据库连接实例
func New(opts ...Option) (*gorm.DB, error) {
	opt := setOption(opts...)
	if len(opt.dbConfigs) != 1 {
		return nil, errors.New("此方法只能实例化一个数据库连接实例")
	}

	return newConnect(&opt.dbConfigs[0], opt)
}

// NewMulti 批量初始化数据库连接实例，成功会返回多个以数据库名为key，*gorm.DB为value的集合
func NewMulti(opts ...Option) (map[string]*gorm.DB, error) {
	opt := setOption(opts...)

	if len(opt.dbConfigs) < 1 {
		return nil, errors.New("初始化的数据库配置个数不能为0")
	}

	dbs := make(map[string]*gorm.DB)
	for _, cfg := range opt.dbConfigs {
		conn, err := newConnect(&cfg, opt)
		if err != nil {
			logger.Logger.Error("New Mysql Conn Err:", zap.Error(err))
			return nil, err
		}

		dbs[cfg.DBName] = conn
	}

	return dbs, nil
}

// setOption 设置option
func setOption(opts ...Option) *option {
	opt := &option{
		maxIdleConn:     defaultMaxIdleConn,
		maxOpenConn:     defaultMaxOpenConn,
		connMaxLifetime: defaultConnMaxLifetime,
	}

	for _, f := range opts {
		f(opt)
	}

	return opt
}

// newConnect 初始化一个数据库连接
// TODO 连接问题
func newConnect(cfg *model.PolarConfig, opt *option) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=2s&readTimeout=6s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &opt.gormConfig)
	if err != nil {
		return nil, err
	}

	// 初始化连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConn 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(opt.maxIdleConn)
	// SetMaxOpenConn 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(opt.maxOpenConn)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(opt.connMaxLifetime)

	return db, nil
}
