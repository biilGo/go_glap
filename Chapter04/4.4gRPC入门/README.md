# go_glap
Go-language-advanced-programming

## gRPC入门
基于http/2链接提供多个服务,对于移动设备更加友好.

### gRPC技术栈
![ch4-1-grpc-go-stack.jpg](./assets/ch4-1-grpc-go-stack.jpg)

最底层为TCP或unix socket协议,在此之上http/2协议的实现,然后在http/2协议之上又构建了针对gRPC核心库.应用程序通过gRPC插件生产的Stub代码和gRPC核心库通信,也可以直接和gRPC核心库通信

### gRPC流
rpc是远程函数调用,因此每次调用的函数参数和返回值不能太大,否则将严重影响每次调用的响应时间.因此传统的rpc方法调用对于上传和下载较大数据量场景并不合适,同时传统RPC模式也不适用于对时间不确定的订阅和发布模式.为此,gRPC框架针对服务端和客户端分别提供了流特性.