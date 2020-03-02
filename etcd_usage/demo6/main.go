package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config          clientv3.Config
		client          *clientv3.Client
		err             error
		kv              clientv3.KV
		watcher         clientv3.Watcher
		getResp         *clientv3.GetResponse
		watcherRevision int64
		watcherChan     <-chan clientv3.WatchResponse
		watcherResp     clientv3.WatchResponse
		event           *clientv3.Event
		cxt             context.Context
		cancelFunc      context.CancelFunc
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

	go func() {
		for {
			kv.Put(context.TODO(), "/new/json/rose", "玫瑰")

			kv.Delete(context.TODO(), "/new/json/rose")

			time.Sleep(time.Second)
		}
	}()

	// 先get到当前的值 并监听后期变化
	if getResp, err = kv.Get(context.TODO(), "/new/json/rose"); err != nil {
		fmt.Println(err)
		return
	}
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值为：", string(getResp.Kvs[0].Value))
	}
	// 监听版本
	watcherRevision = getResp.Header.Revision + 1

	//创建一个watcher
	watcher = clientv3.NewWatcher(client)

	fmt.Println("从该版本向后监听", watcherRevision)

	// 取消监听
	cxt, cancelFunc = context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	watcherChan = watcher.Watch(cxt, "/new/json/rose", clientv3.WithRev(watcherRevision))

	for watcherResp = range watcherChan {
		for _, event = range watcherResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为值：", string(event.Kv.Value), "Version:", event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了：", "Version:", event.Kv.ModRevision)
			}
		}
	}
}
