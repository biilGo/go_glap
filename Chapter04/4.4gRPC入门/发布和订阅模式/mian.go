// docker项目中提供了一个pubsub的极简实现

package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	// pubsub.NewPublisher构造一个发布对象,p.SubscribeTopic()可以通过函数筛选感兴趣的主题进行订阅
	p := pubsub.NewPublisher(100*time.Millisecond, 10)

	golang := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "golang:") {
				return true
			}
		}
		return false
	})
	docker := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "docker:") {
				return true
			}
		}
		return false
	})

	go p.Publish("hi")
	go p.Publish("golang: https://golang.org")
	go p.Publish("docker: https://www.docker.com/")
	time.Sleep(1)

	go func() {
		fmt.Println("golang topic:", <-golang)
	}()

	go func() {
		fmt.Println("docker topic:", <-docker)
	}()

	<-make(chan bool)
}

// 基于grpc和pubsub包,提供一个跨网络的发布和订阅系统
service PubsubService {
	rpc Publish (String) returns (String);
	rpc Subscribe (String) returns (stream String);
}

// Publish是普通的rpc方法,Subscribe是一个单向的流服务.然后grpc插件会为服务端和客户端生成对应的接口:
type PublishServiceServer interface {
	Publish(context.Context, *String) (*String, error)
	Subscribe(*String, PubsubService_SubscribeServer) error
}

type PubsubServiceClient interface {
	Publish(context.Context, *String, ...grpc.CallOption) (*String, error)
	Subscribe(context.Context, *String, ...grpc.CallOption) (
		PubsubService_SubscribeClient, error,
	)
}

type PubsubService_SubscribeServer interface {
	Send(*String) error
	grpc.grpcServerStream
}

// Subscribe是服务端的单向流,因此生成的HelloService_SubscribeServer接口中只有Send方法
// 然后就可以实现发布和订阅服务:
type PubsubService struct {
	pub *pubsub.Publisher
}

func NewPubsubService() *PubsubService {
	return &PubsubsService {
		pub: pubsub.NewPublisher(100*time.Millisecond,10),
	}
}

// 实现发布方法和订阅方法:
func (p *PubsubService) Publish(
	ctx context.Context, arg *String,
) (*String, error) {
	p.pub.Publish(arg.GetValue())
	return &String{},nil
}

func (p *PubsubService) Subscribe(
	arg *String, stream PubsubService_SubscribeServer,
) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if string.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		if err := stream.Send(&String{Value: v.(string)}); err != nil {
			return err
		}
	}

	return nil
}

// 从客户端向服务器发布信息
func main()  {
	conn, err := gtpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewPubsubServiceClient(conn)

	_, err := client.Publish(
		context.Background(), &String{Value: "golang: hello Go"},
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish(
		context.Background(), &String{Value: "docker: hello Docker"},
	)

	if err != nil {
		log.Fatal(err)
	}
}

// 可以在另一个客户端进行订阅信息：
func main()  {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewPubsubServiceClient(conn)
	stream, err := client.SubscribeTopic(
		context.Background(), &String{Value: "golang:"},
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		replpy, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}