## Introduction

SOCKS5 是一个代理协议，旨在为位于 Intranet 防火墙后的用户提供访问 Internet 的代理服务（Intranet，你没听错，这是个有一定年头的协议，其 RFC 提案的时间比 HTTP 1.0 还要早两个月）。

> 代理
> 根据 HTTP 1.1 的定义，proxy 是：
> An intermediary program which acts as both a server and a client for the purpose of making requests on behalf of other
> clients. Requests are serviced internally or by passing them on, with possible translation, to other servers.
>
> 代理就是中间人，一人分饰两角：客户端眼中的目标服务器，目标服务器眼中的客户端——这意味着他必须同时满足C/S 双方的规范。再细分，如果只是简单的 pipe C/S 两端数据，那他就是个“透明代理”；一旦他对请求或响应进行了修改，那就是“非透明代理”。

但其实，SOCKS5 协议并不负责代理服务器的数据传输环节，此协议只是在C/S两端真实交互之间，建立起一条从客户端到代理服务器的授信连接。 refer links:

### 协议流程

从流程上来说，SOCKS5 是一个C/S 交互的协议，交互大概分为这么几步：

1. 客户端发送认证协商
2. 代理服务器就认证协商进行回复（如拒绝则本次会话结束）
    1. 如需GSSAPI或用户名/密码认证，客户端发送认证信息
    2. 代理服务器就对应项进行鉴权，并进行回复或拒绝
3. 客户端发送希望连接的目标信息
4. 代理服务器就连接信息进行确认或拒绝
5. 非协议内容】：代理服务器连接目标并 pipe 到客户

### refer links

- https://www.rfc-editor.org/rfc/rfc1928.html
- https://jiajunhuang.com/articles/2019_06_06-socks5.md.html
- http://www.moye.me/2017/08/03/analyze-socks5-protocol/
- https://segmentfault.com/a/1190000038247560