## Introduction

本示例主要是演示使用Protobuf协议去实现Go下的RPC,而非gRPC。其中涉及到自定义的Protobuf工具集的go插件,其是利用原有的protoc-gen-go插件进行 自定义改造，使生成的go代码非基于gRPC，而是RPC方式。

> 对于Go语言的protoc-gen-go插件来说，里面又实现了一层静态插件系统。比如protoc-gen-go内置 了一个gRPC插件，用户可以通过
> --go_out=plugins=grpc参数来生成gRPC相关代码，否则只会针对message生成相关代码,而不会针对service生成server和client对应代码。

- hello-server: 使用了根据hello.proto文件生成的Protobuf协议的接口(不包含server和client端)
- protoc-gen-go-netrpc: Protobuf 工具集的插件，能够根据.proto文件生成使用go的RPC包，和Protobuf协议的RPC代码，包含接口部分,server和client