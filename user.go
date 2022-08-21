package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	Conn   net.Conn
	server *Server
}

// 用户上线
func (this *User) online() {
	//但是此时无法直接使用server的map，因此在user类里面加上server
	this.server.rwlock.Lock()
	this.server.UserMap[this.Name] = this
	this.server.rwlock.Unlock()
	this.server.broadcast(this, "上线")
}

// 用户下线
func (this *User) offline() {
	this.server.rwlock.Lock()
	delete(this.server.UserMap, this.Name)
	this.server.rwlock.Unlock()
	this.server.broadcast(this, "下线")
}

func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg))
}

// 用户发消息
func (this *User) DoMessage(msg string) {
	//查询功能 查询当前在线人
	if msg == "who" {
		for _, user := range this.server.UserMap {
			onlinemsg := "[" + user.Name + "]" + "在线" + "\n"
			this.SendMsg(onlinemsg)
		}
	} else if len(msg) > 7 && msg[:6] == "rename" {
		//改名功能形式为:"rename,newname"
		newname := strings.Split(msg, ",")[1]
		this.server.rwlock.Lock()
		delete(this.server.UserMap, this.Name)
		this.server.UserMap[newname] = this
		this.server.rwlock.Unlock()
		this.SendMsg("您已经成功修改名字,新名字" + newname + "\n")
	} else if len(msg) > 4 && msg[:3] == "to|" {
		this.prichat(msg)
	} else {
		this.server.MSG <- msg
	}

}

// 私聊
func (this *User) prichat(msg string) {
	if strings.Split(msg, "|")[1] == " " {
		this.DoMessage("您输入的方式有误，请按照'to|zhangsan|你好'这种格式输入")
		return
	} else if strings.Split(msg, "|")[2] == " " {
		this.DoMessage("不可发送空消息")
		return
	} else {
		Tousername := strings.Split(msg, "|")[1]
		content := strings.Split(msg, "|")[2]
		Touser := this.server.UserMap[Tousername]
		Touser.SendMsg(this.Name + "给您发的消息:" + content)
	}
}

// 创建一个用户的接口
func NewUser(conn net.Conn, conserver *Server) *User {
	UserAddr := conn.RemoteAddr().String()
	//第一次这次写的LocalAddr 导致ip与port都一样，
	user := &User{
		Name:   UserAddr,
		Addr:   UserAddr,
		C:      make(chan string),
		Conn:   conn,
		server: conserver,
	}

	go user.ListenuMsg()

	return user
}

// 用户的监听方法
// 一听到消息就写到用户自己的管道里
func (this *User) ListenuMsg() {
	for {
		msg := <-this.C

		this.Conn.Write([]byte(msg + "\n"))
	}
}
