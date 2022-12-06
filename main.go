package main

import (
	"colly_go_demo/demos/llss"
	"colly_go_demo/demos/xunacg"
	"fmt"
)

func main() {

	var demoVal int
	fmt.Println("请输入要运行的采集器")
	fmt.Println("1 琉璃神社")
	fmt.Println("2 xunacg（签到）")
	fmt.Scanln(&demoVal) // 接收

	if demoVal == 1 {
		llss.Run()
	} else {
		xunacg.SignIn()
	}
}
