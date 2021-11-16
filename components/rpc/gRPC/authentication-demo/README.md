## Introduction

gRPC建立在HTTP/2协议之上，对TLS提供了很好的支持。本模块外的gRPC的服务都没有提供证书支持，因此客户端在链接服务器中通过grpc.WithInsecure()选 项跳过了对服务器证书的验证

> Notes: SSL 是指安全套接字层，简而言之，它是一项标准技术，可确保互联网连接安全，保护两个系统之间发送的任何敏感数据，防止网络犯罪分子读取和修改
> 任何传输信息，包括个人资料。两个系统可能是指服务器和客户端（例如，浏览器和购物网站），或两个服务器之间（例如，含个人身份信息或工资单信息的应用程
> 序）。 此举可确保在用户和站点之间，或两个系统之间传输的数据无法被读取。它使用加密算法打乱传输中的数据，防止数据通过连接传输时被黑客读取。这里所
> 说的数据是指任何敏感或个人信息，例如信用卡号和其他财务信息、个人姓名和住址等。
>
> TLS（传输层安全）是更为安全的升级版 SSL。由于 SSL 这一术语更为
> 常用，因此我们仍然将我们的安全证书称作 SSL。但当您从DigiCert购买 SSL 时，您真正购买的是最新的 TLS 证书，有 ECC、RSA 或 DSA 三种加密方式
> 可以选择。

### Private Key And Certificate

可以用以下命令为服务器和客户端分别生成私钥和证书：

````shell
# 生成私钥
openssl genrsa -out server.key 2048
# 生成证书或公钥
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" \
    -key server.key -out server.crt

# 生成私钥
openssl genrsa -out client.key 2048
# 生成证书或公钥
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io" \
    -key client.key -out client.crt
````

以上命令将生成server.key、server.crt、client.key和client.crt四个文件。其中以.key为后缀名的是私钥文件，需要妥善保管。以.crt为后缀名是证书
文件，也可以简单理解为公钥文件，并不需要秘密保存。在subj参数中的/CN=server.grpc.io表示服务器的名字为server.grpc.io，在验证服务器的证书时需要用到该信息。

### 1-On-web

一个简单使用TSL协议的web demo:

1. 启动一个web服务，其使用server.crt证书和server.key私钥用于TSL加密保护信息安全。 参见main.go的startServer()函数。

2. 启动一个client客户端链接其使用server.crt证书和"server.grpc.io"构建TSL客户端。参见main.go的doClientWork()函数。

````go
func main() {
creds, err := credentials.NewClientTLSFromFile(
"server.crt", "server.grpc.io",
)

if err != nil { log.Fatal(err)
}

conn, err := grpc.Dial("localhost:5000",
grpc.WithTransportCredentials(creds),
)
if err != nil {
log.Fatal(err)
}
defer conn.Close()

...

}
````

其中redentials.NewClientTLSFromFile是构造客户端用的证书对象，第一个参数是服务器的证书文件，第二个参数是签发证书的服务器的名字。然后通过grpc.WithTransportCredentials(
creds)将证书对象转为参数选项传人grpc.Dial函数。

### TLS-CA

以上(On-web中)这种方式，需要提前将服务器的证书告知客户端，这样客户端在链接服务器时才能进行对服务器证书认证。在复杂的网络环境中，服务器证书的传输本身也是一个
非常危险的问题。如果在中间某个环节，服务器证书被监听或替换那么对服务器的认证也将不再可靠。 为了避免证书的传递过程中被篡改，可以通过一个安全可靠的根证书分别对服务器和客户端的证书进行签名。这样客户端或服务器在收到对方的证书后可以通过根证书
进行验证证书的有效性。

- 根证书的生成方式和自签名证书的生成方式类似:

```shell
# 生成根证书的私钥
openssl genrsa -out ca.key 2048

# 生成根证书或公钥
openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=gobook/CN=github.com" \
    -key ca.key -out ca.crt
    
# 重新对服务器端证书进行签名
## 生成服务端证书签名请求文件server.csr(依然使用上面已生成的server.key私钥)
openssl req -new \
    -subj "/C=GB/L=China/O=server/CN=server.io" \
    -key server.key \
    -out server.csr
## 生成服务端证书文件或公钥，生成时引入签名请求文件server.csr，在证书签名完成之后可以删除.csr文件
openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in server.csr \
    -out server.crt    
```

客户端就可以基于CA证书对服务器进行证书验证:

```go
func doClientWork() {
certificate, err := tls.LoadX509KeyPair(client_crt, client_key)
if err != nil {
log.Panicf("could not load client key pair: %s", err)
}

certPool := x509.NewCertPool()
ca, err := ioutil.ReadFile(ca)
if err != nil {
log.Panicf("could not read ca certificate: %s", err)
}
if ok := certPool.AppendCertsFromPEM(ca); !ok {
log.Panic("failed to append ca certs")
}

creds := credentials.NewTLS(&tls.Config{
InsecureSkipVerify: false,         // NOTE: this is required!
ServerName:         tlsServerName, // NOTE: this is required!
Certificates:       []tls.Certificate{certificate},
RootCAs:            certPool,
})

conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(creds))
if err != nil {
log.Fatal(err)
}
defer conn.Close()
...
}
```

在新的客户端代码中，我们不再直接依赖服务器端证书文件。在credentials.NewTLS函数调用中，客户端通过引入一个CA根证书和服务器的名字来实现对服务器进行验证。客户端在链接服务器时会首先请求服务器的证书，然后使用CA根证书对收到的服务器端证书进行验证。

如果客户端的证书也采用CA根证书签名的话，服务器端也可以对客户端进行证书认证。我们用CA根证书对客户端证书签名:

```shell
openssl req -new \
    -subj "/C=GB/L=China/O=client/CN=client.io" \
    -key client.key \
    -out client.csr
openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in client.csr \
    -out client.crt
```

- 因为引入了CA根证书签名，在启动服务器时同样要配置根证书，代码详情可参考main.go的startServer()函数。

### token

前面的基于证书的认证是针对每个gRPC链接的认证。gRPC还为每个gRPC方法调用提供了认证支持，这样就基于用户Token对不同的方法访问进行权限管理。
要实现对每个gRPC方法进行认证，需要实现grpc.PerRPCCredentials接口：

```go
type PerRPCCredentials interface {
GetRequestMetadata(ctx context.Context, uri ...string) (
map[string]string, error,
)
RequireTransportSecurity() bool
}
```

在GetRequestMetadata方法中返回认证需要的必要信息。RequireTransportSecurity方法表示是否要求底层使用安全链接。在真实的环境中建议必须要求底层启用安全的链接，否则认证信息有泄露和被篡改的风险。

> N0otes: 本示例未启用安全链接。

### panic-and-log