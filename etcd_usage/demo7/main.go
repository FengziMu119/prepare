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
		putOp  clientv3.Op
		getOp  clientv3.Op
		opResp clientv3.OpResponse
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

	putOp = clientv3.OpPut("/new/json/lom", "222")
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入的Revision:", opResp.Put().Header.Revision)

	getOp = clientv3.OpGet("/new/json/lom")

	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("数据Revision：", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据Value：", string(opResp.Get().Kvs[0].Value))
}
