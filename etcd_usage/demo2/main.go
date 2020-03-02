package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putRes *clientv3.PutResponse
		//getRes *clientv3.GetResponse
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.31.233:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 创建客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.KV(client)
	if putRes, err = kv.Put(context.TODO(), "/new/json/tom", "hello", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putRes.Header.Revision)
		if putRes.PrevKv != nil {
			fmt.Println(string(putRes.PrevKv.Value))
		}
	}
}
