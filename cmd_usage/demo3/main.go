package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	var (
		cmd        *exec.Cmd
		ctx        context.Context
		cancelFunc context.CancelFunc
		resultChan chan *result
		res        *result
	)
	// 创建一个结果队列
	resultChan = make(chan *result, 1000)
	//执行一个cmd，让他在一个协程里去执行，让他执行2秒 sleep 2 ;echo hello
	//context :chan byte
	//cancelFunc : close(chan byte)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "C:\\cygwin64\\bin\\bash", "-c", "sleep 2;echo hello")
		//执行任务捕获输出
		output, err = cmd.CombinedOutput()
		//把任务输出结果 传给main协程
		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()
	// select{case <-ctx.Done():}
	// kill pid 进程ID 杀死子进程
	// 继续往下走
	time.Sleep(time.Second)
	// 取消上下文
	cancelFunc()

	//在main协程里，等待子协程退出，并打印任务结果
	res = <-resultChan
	fmt.Println(res.err, string(res.output))
}
