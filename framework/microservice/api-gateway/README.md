## Introduction

微服务网关是一个处于应用程序或服务之前的系统，这样微服务就会被微服务网关给保护起来，对所有的调用者透明。

### What is API Gateway?

微服务网关可以这样简单理解:
> 服务网关 = 路由转发 + 过滤器

1. 路由转发：接收一切外界请求，转发到后端的微服务上去.
2. 过滤器：在服务网关中可以完成一系列的横切功能，例如权限校验、限流以及监控等，这些都可以通过过滤器完成（其实路由转发也是通过过滤器实现的.

### Why we need the API Gateway?

客户端可能是你的应用程序的前端，形式是网页或应用程序，也可能是你的组织内部需要与你的应用程序交互的其他内部服务，或者是
第三方客户端应用程序和网站。像API代理一样，网关接收传入的请求，并将其引导到系统的相关部分，然后将响应转发回客户端。但API
网关不仅仅是一个简单的反向代理服务，它提供了一个统一的接口，并提供了安全、负载均衡、请求和响应转换、监控和追踪等功能。

1. 请求接入: 管理所有的接入请求，是所有API接口的请求入口。作为企业系统边界，隔离外网系统与内网系统。
2. 解耦: 通过解耦，使得微服务系统的各方能够独立，自由，高效，灵活的调整，而不用担心给其他方面带来影响。
3. 拦截策略: 提供了一个扩展点，方便通过扩展机制对请求进行一系列的加工和处理。可以提供统一的安全，路由和流控等公共服务组件。
4. 统一管理: 可以通过统一的监控工具，配置管理等基础设施。

### API Gateway Components

#### Zuul

Zuul 是由 Netflix 所开源的组件，基于JAVA技术栈开发的。

Zuul网关的使用热度非常高，并且也集成到了 Spring Cloud 全家桶中了，使用起来非常方便。

![avatar](https://gitee.com/xuzimian/Image/raw/master/Spring/SpringCloud/zuul.png)

请求过程

- 首先将请求给zuulservlet处理，zuulservlet中有一个zuulRunner对象，该对象中初始化了RequestContext：作为存储整个请求的一些数据，并被所有的zuulfilter共享。
- zuulRunner中还有 FilterProcessor，FilterProcessor作为执行所有的zuulfilter的管理器。
- FilterProcessor从filterloader 中获取zuulfilter，而zuulfilter是被filterFileManager所加载，并支持groovy热加载，采用了轮询的方式热加载。
- 有了这些filter之后，zuulservelet首先执行的Pre类型的过滤器，再执行route类型的过滤器，最后执行的是post 类型的过滤器，如果在执行这些过滤器有错误的时候则会执行error类型的过滤器。
- 执行完这些过滤器，最终将请求的结果返回给客户端。
- 主要特色是，这些过滤器可以动态插拔，就是如果需要增加减少过滤器，可以不用重启，直接生效。
- 原理就是：通过一个db维护过滤器（上图蓝色部分），如果增加过滤器，就将新过滤器编译完成后push到db中，有线程会定期扫描db，发现新的过滤器后，会上传到网关的相应文件目录下，并通知过滤器loader进行加载相应的过滤器。

#### Tyk

Tyk是一个基于GO编写的，轻量级、快速可伸缩的开源的API网关。支持配额和速度限制，支持认证和数据分析，支持多用户多组织，提供全 RESTful API

#### Kong

Kong是基于OpenResty技术栈的开源网关服务，因此其也是基于Nginx实现的。 Kong可以做到高性能、插件自定义、集群以及易于使用的Restful API管理。

#### Traefik

Traefik 是一个现代 HTTP 反向代理和负载均衡器，可以轻松部署微服务。 Traeffik 可以与现有的组件（Docker、Swarm，Kubernetes，Marathon，Consul，Etcd，…）集成，并自动动态配置

#### Ambassador

Ambassador 是一个开源的微服务 API 网关，建立在 Envoy 代理之上，为用户的多个团队快速发布，监控和更新提供支持 支持处理 Kubernetes ingress controller 和负载均衡等功能，可以与 Istio
无缝集成


#### API Gateway components vs

![avatar](https://gitee.com/xuzimian/Image/raw/master/Spring/SpringCloud/API-Gateway-vs-1.png)
![avatar](https://gitee.com/xuzimian/Image/raw/master/Spring/SpringCloud/API-Gateway-vs-2.png)
![avatar](https://gitee.com/xuzimian/Image/raw/master/Spring/SpringCloud/API-Gateway-vs-3.png)




