package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))

}

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitMySQL() {
	//自定义日志模板 打印SQL
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold: time.Second, //慢SQL阈值
	//		LogLevel:      logger.Info, //级别
	//		Colorful:      true,        //彩色
	//	},
	//)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	fmt.Println(" MySQL inited 。。。。")
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

}

const (
	PublishKey = "websocket"
)

// 发布消息
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	err = Red.Publish(ctx, channel, msg).Err()
	fmt.Println("Publish 。。。。", msg)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// 订阅消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe 。。。。", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return msg.Payload, err
}
