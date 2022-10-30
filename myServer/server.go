package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	ip   string
	port int
	// 在线用户列表
	onlineMap map[string]*User
	// 加个锁
	mapLock sync.RWMutex
	// 消息广播
	channelMsg chan string
}

// 新建Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		ip:         ip,
		port:       port,
		onlineMap:  make(map[string]*User),
		channelMsg: make(chan string),
	}
	return server
}

// get
func (s *Server) GetIp() string {
	return s.ip
}

func (s *Server) GetPort() int {
	return s.port
}

func (s *Server) GetOnlineMap() map[string]*User {
	return s.onlineMap
}

func (s *Server) GetChannelMsg() chan string {
	return s.channelMsg
}

// set
func (s *Server) SetIp(ip string) {
	s.ip = ip
}

func (s *Server) SetPort(port int) {
	s.port = port
}

func (s *Server) SetOnlineMap(onlineMap map[string]*User) {
	s.onlineMap = onlineMap
}

func (s *Server) SetChannelMsg(msg chan string) {
	s.channelMsg = msg
}

// 监听广播信息
func (s *Server) ListenMsger() {
	for {
		msg := <-s.channelMsg
		// 将Msg发送给全部在线用户
		s.mapLock.Lock()
		for _, user := range s.onlineMap {
			user.channelMsg <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息
func (s *Server) BroadCast(u *User, msg string) {
	sendMsg := "[" + u.addr + "]" + u.name + ":" + msg
	s.channelMsg <- sendMsg
}

// 链接业务
func (s *Server) Handler(conn net.Conn) {
	// fmt.Println("===链接建立成功===")
	// 1. 获取用户信息
	u := NewUser(conn, s)
	// 2. 用户上线
	u.Online()
	// 监听用户是否活跃
	isLive := make(chan bool)

	// 持续接收用户客户端信息
	go func() {
		bufferMsg := make([]byte, 4096) // 字节数组，上限4k
		for {
			n, err := conn.Read(bufferMsg)
			if n == 0 {
				u.Offline() // 用户下线
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Error:", err)
			}
			msg := string(bufferMsg[:n-1]) // 转换为字符串，去除换行符
			// 消息广播
			u.DoMsg(msg)
			isLive <- true // 活跃
		}
	}()
	// 当前handler阻塞
	for {
		select {
		case <-isLive:

		case <-time.After(time.Minute * 5): //5分钟不活跃
			u.SendMsg("=====你长时间未活跃,被强制下线!=====\n")
			close(u.channelMsg)
			conn.Close()
			return
		}
	}
}

// 启动服务器方法
func (s *Server) Start() {
	// socket listen
	listener, listen_err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if listen_err != nil {
		fmt.Println("net.Listen Error: ", listen_err)
		return
	}
	defer listener.Close() // close
	go s.ListenMsger()     // 启动监听

	for {
		// accept
		conn, conn_err := listener.Accept()
		if conn_err != nil {
			fmt.Println("listener.Accpet Error: ", conn_err)
			continue
		}
		// do handler
		go s.Handler(conn)
	}
}
