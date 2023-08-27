package main

// 服务端或客户端的单向流是双向流的特例,在HelloService增加一个支持双向流的Channel方法
// 关键字stream指定启用流特性,参数部分是接收客户端参数的流,返回值是返回给客户端的流
service HelloService {
	rpc Hello (String) returns (String);

	rpc Channel (stream String) returns (stream String);
}

// 重新生成代码可以看到接口种新增加的Channel方法的定义:
// 在服务端的Channel方法参数是一个新的HelloService_ChannelServer类型的参数,可以用于和客户端双向通信.客户端的Channel方法返回一个Hello Service_ChannelClient类型的返回值,可以用于和服务端进行双向通信
type HelloServiceServer interface {
	Hello(context.Context, *String) (*String, error)
	Channel(HelloService_ChannelServer) error
}

type HelloServiceClient interface {
	Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (
		*String, error,
	)
	Channel(ctx context.Context, opts ...grpc.CallOption) (
		HelloService_ChannelClient, error,
	)
}

// HelloService_ChannelServer和HelloService_ChannelClient均为接口类型:
// 服务端,客户端的流辅助接口均定义了Send和Recv方法用于流数据的双向通信.
type HelloService_ChannelServer interface {
	Send(*String) error
	Recv() (*String,error)
	grpc.grpcServerStream
}

type HelloService_ChannelClient interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ClientStream
}

func (p *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {
	// 服务端在循环中接收客户端发来的数据
	// 如果遇到io.EOF表示客户端流被关闭,如果函数退出表示服务端流关闭.
	// 生成返回的数据通过流发送给客户端,双向流数据的发送和接收都是完全独立的行为.
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &String{Value:"hello:" + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}

	// 客户端需要先调用Channel方法获取返回的流对象:
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 在客户端我们将发送和接收操作放到两个独立的Goroutine.首先是向服务端发送数据:
	go func() {
		for {
			if err := stream.Send(&String{Value:"hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	// 然后再循环中接收服务端返回的数据:
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
}