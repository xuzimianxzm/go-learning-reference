## Install Project

```shell
 go mod init github.com/xuzimian/blog-web-demo

 go get -u github.com/gin-gonic/gin
 
 go get -u github.com/unknwon/com
 
 go get -u github.com/jinzhu/gorm 
 
 go get -u github.com/go-sql-driver/mysql

 go mod tidy
```

## Knowledge Points

### Standard Library

- fmt：实现了类似 C 语言 printf 和 scanf 的格式化 I/O。格式化动作（‘verb’）源自 C 语言但更简单
- net/http：提供了 HTTP 客户端和服务端的实现

### Gin

- gin.Default()：返回 Gin 的type Engine struct{...}，里面包含RouterGroup，相当于创建一个路由Handlers，可以后期绑定各类的路由规则和函数、中间件等
- router.GET(…){…}：创建不同的 HTTP 方法绑定到Handlers中，也支持 POST、PUT、DELETE、PATCH、OPTIONS、HEAD 等常用的 Restful 方法
- gin.H{…}：就是一个map[string]interface{}
- gin.Context：Context是gin中的上下文，它允许我们在中间件之间传递变量、管理流、验证 JSON 请求、响应 JSON
  请求等，在gin中包含大量Context的方法，例如我们常用的DefaultQuery、Query、DefaultPostForm、PostForm等等


* http.ListenAndServe 和 gin的Engine.Run() 没有本质区别，都是执行的http.ListenAndServe()