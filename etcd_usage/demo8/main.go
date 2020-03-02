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
		kv             clientv3.KV
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
		cancelFunc     context.CancelFunc
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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
	//lease实现锁自动过期
	// op操作
	// txn事务 if else then
	//上锁 （创建租约，自动续约，拿着租约去抢占一个key）
	// 申请一个lease(租约)
	lease = clientv3.NewLease(client)

	//申请一个5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}
	//拿到租约的id
	leaseId = leaseGrantResp.ID

	// 取消自动续约
	ctx, cancelFunc = context.WithCancel(context.TODO())
	// 确保函数退出 关闭
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
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
	// 处理业务
	// if 不存在key then 设置它 else 抢锁失败
	txn = kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision("/new/json/mom"), "=", 0)).
		Then(clientv3.OpPut("/new/json/mom", "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/new/json/mom")) // 否则抢锁失败

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}
	if !txnResp.Succeeded {
		fmt.Println("锁被占用：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	time.Sleep(5 * time.Second)
	// defer 会释放租约
}
