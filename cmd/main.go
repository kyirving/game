package main

import (
	"fmt"
	"game/game"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
)

func main() {
	// start := time.Now()
	// game.Run()

	// elapsed := time.Since(start)
	// fmt.Printf("程序运行时间：%s\n", elapsed)

	c := cron.New()
	c.AddFunc("0 */5 * * * *", func() {
		game.Run()
	})
	c.Start()

	// 创建一个信号通道来接收中断信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞程序执行，直到接收到中断信号
	<-signalChan

	// 停止调度器并执行清理操作
	c.Stop()
	fmt.Println("程序已退出")

}
