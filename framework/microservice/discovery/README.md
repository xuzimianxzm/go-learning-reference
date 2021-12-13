## Introduction

- Consul: 是HashiCorp公司推出的开源工，使用Go语言开发，具有开箱即可部署方便的特点。Consul是分布式的、高可用的、 可横向扩展的用于实现分布式系统的 服务发现与配置满足CP特性。

- Etcd: 是由CoreIS开源，采用Go语言编写的分布式，高可用的Key/Value存储系统，主要用于服务发现和配置共享。相比其它组件，Etcd更为轻量级，部署简单 ，支持http接口。

- Zookeeper: 作为Hadoop 和 Hbase 的重要组件，是一个开源的分布式应用协调服务，采用Java语言编写，致力于为分布式应用提供一致性服务。

## 服务注册与发现组件的对比与选型

| 功能点         | Consul    | Etcd     | Zookeeper | 
| ------------- | ------    | -------- | --------- |
| CAP原理        | CP        | CP       | CP        |
| Key/value存储  | 支持       | 支持      | 支持      |
| 多数据中心      | 支持       | 支持      | 支持      |
| 一致性协议      | Raft      | Raft      | ZAB      |
| 访问协议        | HTTP/DNS  | HTTP/Grpc | RPC客户端 |
| Watch 机制     | 支持       | 支持       | 支持     |
| 安全机制        | ACL/HTTPS | HTTPS     | ACL      |
| 健康检查        | 健康检查    | 长链接    | 连接心跳   |

- 从软件的生态出发，Consul是以服务发现和配置作为主要功能目标，附带提供了Key/Value 存储，相对于Etcd 和 Zookeeper 来讲业务番位较小，更适合于 服务发现与注册。

- Etcd 和 Zookeeper 属于通用的分布式一致性存储系统，被应用于分布式系统的协调工作中，使用番位抽象，具体的业务场景需要开发人员自主实现，如服务注 册和发现，分布式锁等。

- 仅从服务注册于发现组件的需求来看，Consul作为服务发现于注册中心效果更好。如果系统中存在其他分布式一致性协作需求，选择Etcd 和 Zookeeper 反而 能够提供更多的服务支持。

## Go-Kit

Go-kit是一套微服务工具集，用于帮助开发人员解决分布式系统开发中的相关问题，是开发人员更专注于业务逻辑的开发。Go-kit中提供了多种服务注册与发现组件 包括 Consul，Etcd, Zookeeper等。

## Architecture

### 本项目定义了两套服务发现客户端

- discover_client: 自定义服务发现与注册客户端接口层。
- my_discover_client: 手动实现的consul交互细节的 服务发现与注册客户端。
- kit_discover_client: 借用Go-kit 实现consul交互的 服务发现与注册客户端。
- string=service: 基于discover_client提供的服务注册与发现能力搭建的微服务。