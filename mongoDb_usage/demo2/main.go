package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:startTime`
	EndTime   int64 `bson:endTime`
}

type LogRecord struct {
	JobName   string    `bson:jobName` //任务名
	Command   string    `bson:command` //shell命令
	Err       string    `bson:err`     // 脚本错误
	Content   string    `bson:content` //脚本输出
	TimePoint TimePoint `bson:timePoint`
}

func main() {
	var (
		clientOptions *options.ClientOptions
		client        *mongo.Client
		err           error
		dataBase      *mongo.Database
		collection    *mongo.Collection
		record        *LogRecord
		insertRes     *mongo.InsertOneResult
		docid         primitive.ObjectID
	)
	// 建立连接
	clientOptions = options.Client().ApplyURI("mongodb://192.168.31.233:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	// 选择数据库
	dataBase = client.Database("cron")
	// 选择表
	collection = dataBase.Collection("log")

	// 插入一条数据（bosn）
	record = &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}
	//// 插入单条数据
	//if insertRes, err = collection.InsertOne(context.TODO(), record); err != nil {
	//    fmt.Println(err)
	//    return
	//}
	//// _id 默认生成  object 类型
	//docid = insertRes.InsertedID.(primitive.ObjectID)
	//fmt.Println("自增id", docid.Hex())

}
