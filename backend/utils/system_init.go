package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited!!!!")
}

func InitMySQL() {
	//自定義日誌模板 打印SQL語句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢查詢
			LogLevel:      logger.Info, //級別
			Colorful:      true,        //彩色
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("Failed to connect to MySQL:", err)
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Failed to get sql.DB from GORM:", err)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		fmt.Println("Failed to ping MySQL:", err)
	} else {
		fmt.Println("MySQL connect successful!!!!")
	}
}

func InitRedis() {
	ctx := context.Background()
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := Red.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to ping Redis......", err)
	} else {
		fmt.Println("Redis connect successful!!!! => ", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 發佈消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish.......", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 訂閱Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe.......", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe.......", msg.Payload)
	return msg.Payload, err

}
