package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/spf13/viper"
	"os"
	"time"
	"utils/db/config"
)

var Client *redis.Client
var ClusterClient *redis.ClusterClient

func init() {
	if err := config.Init(""); err != nil {
		panic(err)
	}

	addr := viper.GetString("redis.addr")
	pwd := viper.GetString("redis.pwd")
	cluster := viper.GetBool("redis.cluster")

	if cluster {
		if pwd != "" {
			ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:        []string{addr},
				Password:     pwd, // no password set
				WriteTimeout: 2 * time.Minute,
			})
		} else {
			ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:        []string{addr},
				WriteTimeout: 2 * time.Minute,
			})
		}
		if err := ClusterClient.Ping(context.Background()).Err(); err != nil {
			logger.Errorf("connect to redis cluster fail:", err)
			os.Exit(1)
		}
	} else {
		if pwd != "" {
			Client = redis.NewClient(&redis.Options{
				Addr:         addr,
				Password:     pwd, // no password set
				DB:           0,   // use default DB
				WriteTimeout: 2 * time.Minute,
			})
		} else {
			Client = redis.NewClient(&redis.Options{
				Addr:         addr,
				DB:           0, // use default DB
				WriteTimeout: 2 * time.Minute,
			})
		}
		if err := Client.Ping(context.Background()).Err(); err != nil {
			logger.Errorf("connect to redis fail:", err)
			os.Exit(1)
		}
	}

}
