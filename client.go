package main

import (
	"fmt"
	"net"
)

type Client struct {
	ip   string
	port int
	conn net.Conn
	name string
}

func (this *Client) NewClient(conip string, conport int) *Client {
	//创建客户端对象
	client := &Client{
		ip:   conip,
		port: conport,
	}
	//连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s,%d", conip, conport))
	if err != nil {
		fmt.Println("net Dial 失败")
		return nil
	}

	client.conn = conn
	return client
}

func main() {
	newcli := NewServer("127.0.0.1", 8888)
	if newcli == nil {
		fmt.Println(">>>>>>连接服务器失败...")
	} else {
		fmt.Println(">>>>>>连接服务器成功...")
	}

	select {}
}
