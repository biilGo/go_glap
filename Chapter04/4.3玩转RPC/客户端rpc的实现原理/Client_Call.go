package main

import "fmt"

// Go语言的rpc库最简单的使用方式通过Client.Call方法进行同步阻塞调用

func (client *Client) Call {
	serviceMethod string, args interface{},
	reply interface{},
} error {
	call := <- client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

// 通过Client.go方法进行一次异步调用,返回一个表示这次调用的Call结构体,然后等待Call结构体的Done管道返回调用结果.
// 也可以通过Client.go方法异步调用前面的HelloService服务:

func doClientWork(client *rpc.Client) {
	helloCall := client.Go("HelloService.Hello", "hello", new(string),nil)

	// do some thing

	helloCall = <- helloCall.Done
	if err := helloCall.Error; err != nil {
		log.Fatal(err)
	}

	args := helloCall.Args.(string)
	reply := helloCall.Reply.(string)
	fmt.Println(args,reply)
}

// 异步调用命令发出后,一般会执行其他的任务,因此异步调用的输入参数和返回值可以通过返回的Call变量进行获取
// 执行异步调用Client.Go方法实现如下:

func (client *Client) GO(
	serviceMethod string, args interface{},
	reply interface{},
	done chan *Call,
) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.Done = make(chan *Call, 10)

	client.send(call)
	return call
}

// 构造一个表示当前调用的call变量,然后通过client.send将call的完整参数发送到RPC框架.client.send方法调用是线程安全的,因此可以从多个Goroutine同时向同一个rpc链接发送调用指令
// 当调用完成或者发生错误时,将调用call.done方法通知完成:

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// we don't want to block here,It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

// Call.done方法的实现可以得知call.Done管道会将处理后的call返回.