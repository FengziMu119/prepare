package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)
	// 生成cmd
	cmd = exec.Command("C:\\cygwin64\\bin\\bash.exe", "-c", "ls-l")

	//执行命令，捕获了子进程的输出(pipe)
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
}
