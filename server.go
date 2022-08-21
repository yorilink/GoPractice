package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const kickout time.Duration = 10

type Server struct {
	Ip      string
	Port    int
	UserMap map[string]*User
	rwlock  sync.RWMutex

	//广播用户上线的消息
	MSG chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:      ip,
		Port:    port,
		UserMap: make(map[string]*User),
		MSG:     make(chan string),
	}
	return server
}

// 用来处理客户端的逻辑
func (this *Server) Handler(conn net.Conn) {
	//这里接收到客户端的上线消息，创建一个用户
	user := NewUser(conn, this)
	user.online()
	//阻塞
	isAlive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		//创建一个类型为byte, 长度为4096的切片
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//用户下线了未加return，go程已经挂掉了还往里面写东西崩溃
				user.offline()
				return
			}
			if err != nil && io.EOF != err {
				fmt.Println("conn Read err\n", err)
				//未加return
				return
			}

			msg := string(buf[:n-1])
			//取从0到n-1的数值 广播到别的客户端
			user.DoMessage(msg)
			isAlive <- true
		}

	}()
	for {
		select {
		case <-isAlive:
		case <-time.After(time.Second * kickout):
			//user.DoMessage("你已下线\n")
			user.offline()
			//因为是多进程，导致这个地方在往server的msg里写消息之后，下一步如果close运行得更快
			//那么会遍历到UserMap,但是连接已经不存在，程序就崩溃了
			close(user.C)
			//宕机原因 使用已经关闭的连接读数据
			conn.Close()
			//一个连接失效要关闭管道
			return
		}
		//这里应该有个default 如果没写default 会一直等待下去 最后死锁
	}
}

func (this *Server) broadcast(user *User, msg string) {
	SendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.MSG <- SendMsg

}

func (this *Server) ListenMsg() {
	for {
		Amsg := <-this.MSG
		//一旦有用户上线，遍历用户map,对每个user玩家里填入数据
		this.rwlock.Lock()
		for _, cli := range this.UserMap {
			cli.C <- Amsg
		}
		this.rwlock.Unlock()
	}
}

// 启动服务器
func (this *Server) Start() {
	//socket listen
	//监听这个ip和端口
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net,listen err:", err)
	}

	//启动监听
	go this.ListenMsg()
	defer listener.Close()
	//close listen socket
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listen accept err:", err)
			continue
		}

		go this.Handler(conn)
	}

}
