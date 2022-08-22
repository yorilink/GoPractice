package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ip   string
	port int
	conn net.Conn
	name string
	flag int //当前client的模式
}

// 这里如果声明为Clinet 自己的函数则要通过类调用
func NewClient(conip string, conport int) *Client {
	//创建客户端对象
	client := &Client{
		ip:   conip,
		port: conport,
		flag: 999,
	}
	//连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", conip, conport))
	if err != nil {
		fmt.Println("net Dial 失败")
		return nil
	}

	client.conn = conn
	return client
}

func menu(client *Client) bool {
	var flag int
	fmt.Println("输入1进入群发")
	fmt.Println("输入2进入私聊")
	fmt.Println("输入3重命名")
	fmt.Println("输入0退出客户端")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输出正确的数字")
		return false
	}
}

func run(client *Client) {
	for client.flag != 0 {
		for menu(client) != true {
		}
		switch client.flag {
		case 1:
			{
				fmt.Println("群发模式")
				break
			}
		case 2:
			{
				fmt.Println("私聊模式")
				break
			}
		case 3:
			{
				fmt.Println("改名")
				break
			}
		}

	}
}

var serverip string
var serverport int

func init() {
	flag.StringVar(&serverip, "ip", "127.0.0.1", "服务器默认ip:127.0.0.1")
	flag.IntVar(&serverport, "port", 8888, "服务器默认端口号:8888")
}

func main() {
	//命令行解析，init在main函数前运行
	flag.Parse()
	newcli := NewClient(serverip, serverport)
	if newcli == nil {
		fmt.Println(">>>>>>连接服务器失败...")
	} else {
		fmt.Println(">>>>>>连接服务器成功...")
	}

	run(newcli)
}
