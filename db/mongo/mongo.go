package mongo

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/sunshibao/go-utils/db/config"
	"gopkg.in/mgo.v2"
	"os"
)

var session *mgo.Session
var database *mgo.Database
var err error

func initEngine() {
	if err := config.Init(""); err != nil {
		panic(err)
	}
	//user := viper.GetString("mongodb.user")
	//pwd := viper.GetString("mongodb.pwd")
	//host := viper.GetString("mongodb.host")
	var err error

	addr := viper.GetString("mongo.uri")
	user := viper.GetString("mongo.user")
	pwd := viper.GetString("mongo.pwd")
	databases := viper.GetString("mongo.database")

	session, err = mgo.Dial("mongodb://" + addr)
	if err != nil {
		fmt.Printf("mognodb连接出错 %s", err)
		os.Exit(1)
	}

	session.SetMode(mgo.Eventual, true)

	if user != "" && pwd != "" {
		myDB := session.DB("admin") //这里的关键是连接mongodb后，选择admin数据库，然后登录，确保账号密码无误之后，该连接就一直能用了
		//出现server returned error on SASL authentication step: Authentication failed. 这个错也是因为没有在admin数据库下登录
		err = myDB.Login(user, pwd)
		if err != nil {
			fmt.Println("Login-error:", err)
			os.Exit(0)
		}
	}

	database = session.DB(databases)
}

func GetMgoDb() *mgo.Database {
	if database == nil {
		initEngine()
	}
	return database
}
