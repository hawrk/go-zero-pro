package svc

import (
	"algo_assess/assess-rpc-server/internal/config"
	"algo_assess/repo"
	"fmt"
	"github.com/Shopify/sarama"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

type ServiceContext struct {
	Config       config.Config
	SyncProducer sarama.SyncProducer
	// DB sqlx.SqlConn   // 不用框架自带的，比较难用
	//DB          *gorm.DB
	OrderAssessRepo   repo.OrderAssessRepo
	AlgoInfoRepo      repo.AlgoInfoRepo
	UserInfoRepo      repo.AccountInfoRepo
	ProfileRepo       repo.AlgoProfileRepo
	OptimizeRepo      repo.AlgoOptimizeRepo
	OptimizeBaseRepo  repo.AlgoOptimizeBaseRepo
	SummaryRepo       repo.AlgoSummaryRepo
	TimeLineRepo      repo.AlgoTimeLineRepo
	SecurityRepo      repo.SecurityInfoRepo
	AlgoOrderRepo     repo.AlgoOrderRepo
	OrderDetailRepo   repo.OrderDetailRepo
	MarketLevelRepo   repo.MarketLevelRepo
	SHMarketLevelRepo repo.SHMarketLevelRepo
	AccountInfoRepo   repo.AccountInfoRepo
	//RedisClient       *redis.Redis
	ProfileOrigRepo  repo.AlgoProfileOrigRepo  // 算法画像明细表
	TimeLineOrigRepo repo.AlgoTimeLineOrigRepo // 时间线图表
	SummaryOrigRepo  repo.AlgoSummaryOrigRepo  // 汇总表
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化kafka
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Timeout = 5 * time.Second
	kafkaConfig.Producer.Partitioner = sarama.NewManualPartitioner
	p, err := sarama.NewSyncProducer(c.Kafka.Addrs, kafkaConfig)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
	}
	// 初始化Gorm
	level := logger.Default.LogMode(logger.Warn)
	if c.Mysql.EnablePrintSQL == 1 {
		level = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "tb_",
		},
		Logger: level,
	})
	if err != nil {
		fmt.Println("init gorm fail:", err)
		//panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.Mysql.IdleConn)    // 空闲连接池数量
	sqlDB.SetMaxOpenConns(c.Mysql.MaxOpenConn) // 最大连接数量
	sqlDB.SetConnMaxLifetime(time.Hour)        // 设置可复用最大时间

	return &ServiceContext{
		Config:            c,
		SyncProducer:      p,
		OrderAssessRepo:   repo.NewAlgoAssessRepo(db),
		AlgoInfoRepo:      repo.NewAlgoInfoRepo(db),
		UserInfoRepo:      repo.NewAccountInfo(db),
		ProfileRepo:       repo.NewAlgoProfile(db),
		OptimizeRepo:      repo.NewAlgoOptimizeRepo(db),
		OptimizeBaseRepo:  repo.NewAlgoOptimizeBaseRepo(db),
		SummaryRepo:       repo.NewAlgoSummary(db),
		TimeLineRepo:      repo.NewAlgoTimeLine(db),
		SecurityRepo:      repo.NewSecurityInfo(db),
		AlgoOrderRepo:     repo.NewDefaultAlgoOrder(db),
		OrderDetailRepo:   repo.NewOrderDetailRepo(db),
		MarketLevelRepo:   repo.NewMarketLevelRepo(db),
		SHMarketLevelRepo: repo.NewSHMarketLevelRepo(db),
		AccountInfoRepo:   repo.NewAccountInfo(db),
		//RedisClient:       redis.New(c.Redis.Host),
		ProfileOrigRepo:  repo.NewAlgoProfileOrig(db),
		TimeLineOrigRepo: repo.NewAlgoTimeLineOrig(db),
		SummaryOrigRepo:  repo.NewAlgoSummaryOrig(db),
	}
}
