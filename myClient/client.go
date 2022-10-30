package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	serverIp   string
	serverPort int
	name       string
	conn       net.Conn
	num        int
}

func (c *Client) SetName(name string) {
	c.name = name
}
func (c *Client) SetConn(conn net.Conn) {
	c.conn = conn
}

// new Clinet
func NewClinet(ip string, port int) *Client {
	clinet := &Client{
		serverIp:   ip,
		serverPort: port,
		num:        -1,
	}
	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial 出错: ", err)
		return nil
	} else {
		clinet.SetConn(conn)
	}
	return clinet
}

// 处理服务端返回的消息
func (c *Client) DealRes() {
	io.Copy(os.Stdout, c.conn) // 输出即可
}

// 菜单
func (c *Client) Menu() bool {
	fmt.Println("=====请输入您的选择=====")
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更改名字")
	fmt.Println("4. 查询用户")
	fmt.Println("0. 退出")
	fmt.Println("======================")
	var in_num int
	fmt.Scanln(&in_num)
	if in_num >= 0 && in_num <= 4 {
		c.num = in_num
		return true
	} else {
		fmt.Println("=====输入错误!=====")
		return false
	}
}

// 客户端主业务
func (c *Client) Run() {
	for c.num != 0 {

		for !c.Menu() {
		}
		switch c.num {
		case 1:
			c.PublicChat() // 公聊
		case 2:
			c.PrivateChat() // 私聊
		case 3:
			c.Rename() // 重命名
		case 4:
			c.ShowOnline() // 查询在线用户
		}
	}
}

// 更改用户名业务
func (c *Client) Rename() {
	var name string
	fmt.Print("请输入你想要更改的用户名: ")
	fmt.Scanln(&name)
	if name == "exit" {
		fmt.Println("用户名不允许设置为exit!")
		return
	}
	c.name = name
	sendMsg := "rename>" + c.name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("更改用户名失败: ", err)
		return
	}
}

// 公聊业务
func (c *Client) PublicChat() {
	fmt.Print("发送消息(exit退出): ")
	var msg string
	fmt.Scanln(&msg)
	for msg != "exit" {
		if len(msg) != 0 {
			_, err := c.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("发送公聊信息失败: ", err)
				break
			}
		}
		msg = ""
		fmt.Print("发送消息(exit退出): ")
		fmt.Scanln(&msg)
	}
}

// 私聊业务
func (c *Client) PrivateChat() {
	var name string
	var msg string
	fmt.Print("输入私聊对象(exit退出): ")
	fmt.Scanln(&name)

	for name != "exit" {
		fmt.Print("发送信息(exit退出): ")
		fmt.Scanln(&msg)
		for msg != "exit" {
			if len(msg) != 0 {
				sendMsg := "to>" + name + ">" + msg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("发送私聊信息失败: ", err)
					break
				}
			}
			msg = ""
			fmt.Print("发送消息(exit退出): ")
			fmt.Scanln(&msg)
		}
		fmt.Print("输入私聊对象(exit退出): ")
		fmt.Scanln(&name)
	}
}

// 查询在线用户
func (c *Client) ShowOnline() {
	sendMsg := "show online\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("查询在线用户失败: ", err)
		return
	}
}

// 命令行输入参数
var inIp string
var inPort int

// 初始化
func init() {
	// ./client -ip 127.0.0.1 -port 8989
	flag.StringVar(&inIp, "ip", "127.0.0.1", "设置服务端ip")
	flag.IntVar(&inPort, "port", 8989, "设置服务端port")
}

func main() {
	// 命令行解析
	flag.Parse()
	c := NewClinet(inIp, inPort)
	if c != nil {
		go c.DealRes() // 永久阻塞，处理server回复消息
		fmt.Println("=====链接服务端成功=====")
	} else {
		fmt.Println("=====链接服务端失败=====")
		return
	}

	// 启动客户端业务
	c.Run()
}
