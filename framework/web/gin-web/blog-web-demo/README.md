## Install Project

```shell
 go mod init github.com/xuzimian/blog-web-demo

 go get -u github.com/gin-gonic/gin
 
 go get -u github.com/unknwon/com
 
 go get -u github.com/jinzhu/gorm 
 
 go get -u github.com/go-sql-driver/mysql
 
 go get -u github.com/astaxie/beego/validation
 
 go get -u github.com/dgrijalva/jwt-go
 
 go get -u github.com/fvbock/endless

 go mod tidy
```

## Knowledge Points

- Gin：Golang 的一个微框架，性能极佳。
- beego-validation：本节采用的 beego 的表单验证库，中文文档。
- gorm，对开发人员友好的 ORM 框架，英文文档
- com，一个小而美的工具包。

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

### Gorm

gorm 所支持的回调方法：

- 创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
- 更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
- 删除：BeforeDelete、AfterDelete
- 查询：AfterFind

````go
type Article struct {
TagID int `json:"tag_id" gorm:"index"`
Tag   Tag `json:"tag"`
}
````

- gorm:index，用于声明这个字段为索引，如果你使用了自动迁移功能则会有所影响，在不使用则无影响
- Tag 字段，实际是一个嵌套的struct，它利用TagID与Tag模型相互关联，在执行查询的时候，能够达到Article、Tag关联查询的功能。 Article结构体成员是Tag，可以通过Related进行关联查询
- gorm会通过类名+ID 的方式去找到这两个类之间的关联关系

````go
/*
	Preload就是一个预加载器，它会执行两条 SQL，分别是SELECT * FROM blog_articles;和SELECT * FROM blog_tag WHERE id IN (1,2,3,4);
	那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中，会特别方便，并且避免了循环查询
*/
db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
````

### endless

借助 fvbock/endless 来实现 Golang HTTP/HTTPS 服务重新启动的零停机

endless server 监听以下几种信号量：

- syscall.SIGHUP：触发 fork 子进程和重新启动
- syscall.SIGUSR1/syscall.SIGTSTP：被监听，但不会触发任何动作
- syscall.SIGUSR2：触发 hammerTime
- syscall.SIGINT/syscall.SIGTERM：触发服务器关闭（会完成正在运行的请求）
- endless 正正是依靠监听这些信号量，完成管控的一系列动作