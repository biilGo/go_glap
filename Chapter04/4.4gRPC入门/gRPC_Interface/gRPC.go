package main

import (
	"context"
	"fmt"
	"log"
	"net"
)

// gRPC通过context.Context参数,为每个方法调用提供了上下文支持.客户端在调用方法的时候,可以通过可选的grpc.CallOption类型的参数提供额外的上下文信息.

type HelloServiceServer interface {
	Hello(context.Context, *String) (*String, error)
}

type HelloServiceClient interface {
	Hello(context.Context, *String, ...grpc.CallOption) (*String, error)
}

// 基于服务端的HelloServiceServer接口可以重新实现HelloService服务:
type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
) (*String, error) {
	reply := &String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

// gRPC服务的启动流程和标准库的RPC服务启动流程类似:
func main() {
	// 首先通过grpc.NewServer()构造一个gRPC服务对象,然后通过gRPC插件生成的RegisterHelloServiceServer函数注册我们实现的HelloServiceImpl服务.
	grpcServer := grpc.NewServer()
	RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	// 通过在一个监听端口上提供gRPC服务
	grpcServer.Serve(lis)
}

// 通过客户端链接gRPC服务
func main() {
	// grpc.Dial负责和gRPC服务建立链接
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// NewHelloServiceClient函数基于已经建立的连接构造HelloServiceClient对象.返回的client其实是一个HelloServiceClient接口对象,通过接口定义的方法就可以调用服务端对应的gRPC服务提供的方法
	client := NewHelloServiceClient(conn)

	reply, err := client.Hello(context.Background(), &String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply.GetValue())

	// gRPC和标准库的RPC框架有一个区别,gRPC生成的接口并不支持异步调用.不过我们可以在多个groutine之间安全地共享gRPC底层的http/2链接,因此可以通过在另一个goroutine阻塞调用的方式模拟异步调用.
}
