package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestExit(t *testing.T) {
	var (
		mdb = GetMgoDb()
	)
	//选择表 my_collection
	c := mdb.C("info")

	count, err := c.Find(bson.M{}).Count()

	if err != nil {
		panic(err)
	}
	fmt.Println(count)
}
