package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var (
		clientOptions *options.ClientOptions
		client        *mongo.Client
		err           error
		dataBase      *mongo.Database
		collection    *mongo.Collection
	)
	// 建立连接
	clientOptions = options.Client().ApplyURI("mongodb://192.168.31.233:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	// 选择数据库
	dataBase = client.Database("my_db")
	// 选择表
	collection = dataBase.Collection("my_collection")

	collection = collection

}
