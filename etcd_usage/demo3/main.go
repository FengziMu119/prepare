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
		getRes *clientv3.GetResponse
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.31.233:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.KV(client)

	if getRes, err = kv.Get(context.TODO(), "/new/json/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(getRes.Kvs, getRes.Count)
	}

}
