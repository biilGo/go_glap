package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// rpc是基于C/S结构,RPC的服务端对应网络的服务器,rpc的客户端也对应网络客户端.但是对于一些特殊场景,比如在公司内网提供一个rpc服务,但是在外网无法链接到内网的服务器.
// 这种时候我们可以参考反向代理技术,首先从内网主动链接到外网的TCP服务器,然后基于TCP链接向外网提供RPC服务

func main() {
	rpc.Register(new(HelloService))

	for {
		conn, _ := net.Dial("tcp", "localhost:1234")
		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		rpc.ServerConn(conn)
		conn.Close()
	}

	// 反向rpc的内网服务将不再主动提供TCP监听服务,而是首先主动连接到对方的TCP服务器,然后基于每个建立的TCP链接向对方提供的RPC服务

}

// 而RPC客户端则需要在一个公共的地址提供一个TCP服务,用于接受RPC服务器的链接请求:

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	clientChan := make(chan *rpc.Client)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Accept error:", err)
			}

			clientChan <- rpc.NewClient(conn)
		}
	}()

	doClientWork(clientChan)
}

// 客户端执行RPC调用的操作在doClientWork函数完成:
func doClientWork(clientChan <-chan *rpc.Client) {
	client := <-clientChan
	defer client.Close()

	var reply string
	err = client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
