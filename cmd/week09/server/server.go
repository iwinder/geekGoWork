package main

import "github.com/iwinder/geekGoWork/internal/week09"

func main() {
	// // 方案一
	// week09.RunTcpFixLength("127.0.0.1:8033", 1)
	// // 方案二
	// week09.RunTcpFixLength("127.0.0.1:8033", 2)
	// 方案三
	week09.RunTcpFixLength("127.0.0.1:8033", 3)
}
