package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"github.com/sunshibao/go-utils/db/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var dataBase *gorm.DB

func init() {
	if err := config.Init(""); err != nil {
		panic(err)
	}
	dsn := viper.GetString("mysql.user") + ":" +
		viper.GetString("mysql.pwd") + "@tcp(" +
		viper.GetString("mysql.host") + ":" +
		viper.GetString("mysql.port") + ")/" +
		viper.GetString("mysql.database") + "?charset=utf8&parseTime=True&loc=Local"

	conn, err := sql.Open("mysql", dsn)
	conn.SetConnMaxLifetime(30 * time.Minute) // 连接可重用的最长时间
	conn.SetMaxIdleConns(1000)                // 最大空闲连接数
	conn.SetMaxOpenConns(10000)               //最大打开连接数

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: conn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键关系
		SkipDefaultTransaction:                   true, // 跳过默认事务
		PrepareStmt:                              true, // 查询语句解析缓存，提升查询性能
	})

	if err != nil {
		panic(err)
	}

	if viper.GetString("mysql.debug") == "true" {
		db.Logger = db.Logger.LogMode(logger.Info)
	}

	dataBase = db
}

// GetDatabase 获取数据库连接
func GetDatabase() *gorm.DB {
	return dataBase
}

// Ping 检查与数据库的连接是否有效
func Ping(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := dataBase.DB()
	if err != nil {
		return false
	}

	if err := db.PingContext(ctx); err != nil {
		fmt.Println("mysql ping error:", err)
		return false
	}

	return true
}
