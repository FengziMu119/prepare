package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		kv             clientv3.KV
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
	)

	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.31.233:2379"},
		DialTimeout: 5 * time.Second,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	// 申请一个lease(租约)
	lease = clientv3.NewLease(client)

	//申请一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}
	//拿到租约的id
	leaseId = leaseGrantResp.ID
	//自动续租
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已失效了")
					goto END
				} else { //每秒会续租一次，所以就会收到一次应答
					fmt.Println("收到自动续约应答", keepResp.ID)
				}
			}
		}
	END:
	}()
	// 获得kv对象
	kv = clientv3.NewKV(client)
	//Put一个KV，让它与租约关联起来，从而实现10后自动过期
	if putResp, err = kv.Put(context.TODO(), "/new/json/jack", "2222", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功：", putResp.Header.Revision)

	for {
		if getResp, err = kv.Get(context.TODO(), "/new/json/jack"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("过期了")
			break
		}
		fmt.Println("还没过期")
		time.Sleep(2 * time.Second)
	}
}
