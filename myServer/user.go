package main

import (
	"net"
	"strings"
)

type User struct {
	name       string
	addr       string
	channelMsg chan string
	conn       net.Conn
	server     *Server // 属于哪一个服务器
}

// 通过Conn初始化用户
func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		name:       userAddr, // 初始用户名为地址
		addr:       userAddr,
		channelMsg: make(chan string),
		conn:       conn,
		server:     s,
	}
	go user.ListenMsger() // 启动监听
	return user
}

// get方法
func (u *User) GetName() string {
	return u.name
}
func (u *User) GetAddr() string {
	return u.addr
}
func (u *User) GetChannelMsg() chan string {
	return u.channelMsg
}
func (u *User) GetServer() *Server {
	return u.server
}

// set方法
func (u *User) SetName(name string) {
	u.name = name
}
func (u *User) SetAddr(addr string) {
	u.addr = addr
}
func (u *User) SetChannelMsg(msg chan string) {
	u.channelMsg = msg
}
func (u *User) SetServer(s *Server) {
	u.server = s
}

// 监听用户channel，有消息发送给对方客户端
func (u *User) ListenMsger() {
	for {
		msg := <-u.channelMsg            // 从管道获取信息
		u.conn.Write([]byte(msg + "\n")) // 转换为二进制数组
	}
}

// 上线
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.onlineMap[u.GetName()] = u
	u.server.mapLock.Unlock()
	// 广播用户上线信息
	u.server.BroadCast(u, "上线啦！")
}

// 下线
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.onlineMap, u.GetName()) // 从在线集合中删除
	u.server.mapLock.Unlock()
	// 广播用户下线信息
	u.server.BroadCast(u, "下线啦！")
}

// 给当前客户端发消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// 处理消息
func (u *User) DoMsg(msg string) {
	// 查询当前在线用户
	if msg == "show online" {
		u.server.mapLock.Lock()
		for _, user := range u.server.onlineMap {
			onlineMsg := "[" + user.addr + "]" + user.name + ":" + "在线!\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename>" { // 重命名
		// 格式：rename>张三
		new_name := msg[7:]
		if new_name == "exit" {
			u.SendMsg("====用户名不允许设置为exit=====\n")
			return
		}
		// 判断那么是否存在
		_, has_res := u.server.onlineMap[new_name]
		if has_res {
			u.SendMsg("=====当前用户名已存在=====\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.onlineMap, u.name) // 删除旧的用户名
			u.SetName(new_name)
			u.server.onlineMap[u.name] = u
			u.SendMsg("=====更改用户名成功,您当前用户名为" + u.name + "=====\n")
			u.server.mapLock.Unlock()
		}
	} else if len(msg) > 4 && msg[:3] == "to>" { // 私聊格式：to>xxx>msg
		words := strings.Split(msg, ">")
		if len(words) != 3 {
			u.SendMsg("=====发送格式错误: to>xxx>msg=====\n")
			return
		} else {
			toName := words[1] // 发送对象名称
			getUser, ok := u.server.onlineMap[toName]
			if !ok {
				u.SendMsg("=====发送对象用户不存在=====\n")
				return
			} else {
				toMsg := words[2] // 消息
				if toMsg == "" {
					u.SendMsg("=====消息不能为空=====\n")
					return
				} else {
					getUser.SendMsg("Msg By " + u.name + ": " + toMsg)
				}
			}
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}
