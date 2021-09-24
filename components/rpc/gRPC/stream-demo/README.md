## Introduction

- HelloService: Protobuf 接口定义，根据hello.proto文件使用Protobuf工具集生成的go的gRPC代码文件hello.pd.go，该文件分别copy到client和service下
- client: gRPC提供的stream模式客户端
- server: gRPC提供的stream模式服务端