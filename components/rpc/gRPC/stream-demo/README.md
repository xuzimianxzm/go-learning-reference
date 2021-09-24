## Introduction

- HelloService: Protobuf 接口定义，根据hello.proto文件使用Protobuf工具集生成的go的gRPC代码文件hello.pd.go
- client: gRPC提供的stream模式客户端,客户端的Channel方法返回一个HelloService_ChannelClient类型的返回值，可以用于和服务端进行双向通信。
  在客户端我们将发送和接收操作放到两个独立的Goroutine。首先是向服务端发送数据,然后在循环中接收服务端返回的数据。
- server: gRPC提供的stream模式服务端,在服务端的Channel方法参数是一个新的HelloService_ChannelServer类型的参数，可以用于和客户端双向通信。
  服务端在循环中接收客户端发来的数据，如果遇到io.EOF表示客户端流被关闭，如果函数退出表示服务端流关闭。生成返回的数据通过流发送给客户端，双向流数据
  的发送和接收都是完全独立的行为。需要注意的是，发送和接收的操作并不需要一一对应，用户可以根据真实场景进行组织代码。

