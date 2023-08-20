package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type HelloService struct {
	conn net.conn
}

// 为每个链接启动独立的RPC服务:
func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go func() {
			defer conn.Close()
			p := rpc.NewServer()
			p.Register(&HelloService{conn: conn})
			p.ServeConn(conn)
		}()
	}
}

// Hello方法中就可以根据conn成员识别不同链接的RPC调用:
func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request + ", from" + p.conn.RemoteAddr().String()
	return nil
}

// 基于上下文信息,我们可以方便地为rpc服务增加简单的登录状态的验证:
type HelloService struct {
	conn    net.Conn
	isLogin bool
}

func (p *HelloService) Login(request string, reply *string) error {
	if request != "user:passowrd" {
		return fmt.Errorf("auth failed")
	}
	log.Println("login ok")
	p.isLogin = true
	return nil
}

func (p *HelloService) Hello(request string, reply *string) error {
	if !p.isLogin {
		return fmt.Errorf("please login")
	}
	*reply = "hello:" + request + ", from" + p.conn.RemoteAddr().String()
	return nil
}
