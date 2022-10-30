/*
 *程序主入口
 */

package main

func main() {
	s := NewServer("127.0.0.1", 8989) // 创建一个server
	s.Start()                         // 启动服务
}
