// 基于RPC接口规范编写真实的服务端代码

package main

import (
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

// 重构HelloService服务，第一步需要明确服务的名字和接口:

const HelloServiceName = "pth/to/pkg.HelloService"

type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

type HelloServiceClient struct {
	*rpc.Client
}

var _HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(HelloServiceName+".Hello", request, reply)
}

func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

func main() {
	RegisterHelloService(new(HelloService))

	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatal("accept error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}
