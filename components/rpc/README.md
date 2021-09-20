## Introduction

RPC 是远程过程调用的简称，是分布式系统中不同节点间流行的通信方式。Go 语言的标准库也提供了一个简单的 RPC 实现，Go 语言的 RPC 包的路径为 net/rpc，也就是放在了 net 包目录下面。

### gRPC 是什么

![avatar](https://gitee.com/xuzimian/Image/raw/master/RPC/gRPC.png)

gRPC 是可以跨语言开发的
在 gRPC 客户端可以直接调用不同服务器上的远程程序，使用姿势看起来就像调用本地过程调用一样，很容易去构建分布式应用和服务。客户端和服务端可以分别使用 gRPC 支持的不同语言实现

基于 HTTP2 标准设计，比其他框架更优的地方有:

- 支持长连接，双向流、头部压缩、多复用请求等
- 节省带宽、降低 TCP 链接次数、节省 CPU 使用和延长电池寿命
- 提高了云端服务和 Web 应用的性能
- 客户端和服务端交互透明
- gRPC 默认使用 protobuf 来对数据序列化

### protobuf

Protobuf 是 Protocol Buffers 的简称，它是 Google 公司开发的一种数据描述语言，并于 2008 年对外开源。Protobuf 刚开源时的定位类似于 XML、JSON 等数据描述语言，通过附带工具生成代码并实现将结构化数据序列化的功能。但是我们更关注的是 Protobuf 作为接口规范的描述语言，可以作为设计安全的跨语言 PRC 接口的基础工具。

简单介绍 protobuf 的结构定义包含的 3 个关键字

- 以.proto 做为后缀，除结构定义外的语句以分号结尾
- 结构定义可以包含：message、service、enum，三个关键字
- rpc 方法定义结尾的分号可有可无

Message 命名采用驼峰命名方式，字段是小写加下划线

```protobuf
message ServerRequest {
      required string my_name = 1;
}
```

Enums 类型名采用驼峰命名方式，字段命名采用大写字母加下划线

```protobuf
enum MyNum {
VALUE1 = 1;
VALUE2 = 2;
}
```

Service 与 rpc 方法名统一采用驼峰式命名

```protobuf
service Love {
// 定义 Confession 方法
rpc MyConfession(Request) returns (Response) {}
}
```

在 proto 文件中使用 package 关键字声明包名，默认转换成 go 中的包名与此一致，可以自定义包名，修改 go_package 即可：

```protobuf
syntax = "proto3"; // proto版本

package pb; // 指定包名，默认go中包名也是这个

// 定义Love服务
service Love {
  // 定义Confession方法
  rpc Confession(Request) returns (Response) {}
}

// 请求
message Request {
  string name = 1;
}

// 响应
message Response {
  string result = 1;
}
```

Protobuf 核心的工具集是 C++语言开发的，可以使用 proto 工具将对应的.proto 文件编译成 python 等任何其支持的语言，在官方的 protoc 编译器中并不支持 Go，需要安装相应的插件。
安装好 protoc 环境之后，进入到 proto 文件的目录下，执行如下命令，将 proto 文件编译成 pb.go 文件

```sh
 protoc --go_out=plugins=grpc:. test.proto
```
