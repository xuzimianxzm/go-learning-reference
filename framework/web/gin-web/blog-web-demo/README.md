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
 
 go get -u github.com/swaggo/swag/cmd/swag@v1.6.5
 go get -u github.com/swaggo/gin-swagger@v1.2.0 
 go get -u github.com/swaggo/files
 go get -u github.com/alecthomas/template

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

GORM 本身是由回调驱动的，所以我们可以根据需要完全定制 GORM，以此达到我们的目的，如下：

- 注册一个新的回调
- 删除现有的回调
- 替换现有的回调
- 注册回调的顺序

````go
db.Callback().Delete().Replace("gorm:delete", deleteCallback)
````

### endless

借助 fvbock/endless 来实现 Golang HTTP/HTTPS 服务重新启动的零停机

endless server 监听以下几种信号量：

- syscall.SIGHUP：触发 fork 子进程和重新启动
- syscall.SIGUSR1/syscall.SIGTSTP：被监听，但不会触发任何动作
- syscall.SIGUSR2：触发 hammerTime
- syscall.SIGINT/syscall.SIGTERM：触发服务器关闭（会完成正在运行的请求）
- endless 正正是依靠监听这些信号量，完成管控的一系列动作

### Swagger

在项目根目录中，执行初始化命令

````shell
swag init
````

完毕后会在项目根目录下生成docs

```
docs/
├── docs.go
└── swagger
    ├── swagger.json
    └── swagger.yaml
```

- 编写 API 注释，Swagger 中需要将相应的注释或注解编写到方法上，再利用生成器自动生成说明文件。

gin-swagger 给出的范例：

````go
// @Summary Add a new pet to the store
// @Description get string by ID
// @Accept  json
// @Produce  json
// @Param   some_id     path    int     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} web.APIError "We need ID!!"
// @Failure 404 {object} web.APIError "Can not find ID"
// @Router /testapi/get-string-by-int/{some_id} [get]
````

### Docker

构建 Scratch 镜像 Scratch 镜像，简洁、小巧，基本是个空镜像

mysql docker image:

````shell
docker pull mysql
docker run --link mysql:mysql -p 8000:8000 gin-blog-docker-scratch
docker run --link mysql:mysql -p 8000:8000 blog-web-demo
````

### Makefile

Make 是一个构建自动化工具，会在当前目录下寻找 Makefile 或 makefile 文件。如果存在，会依据 Makefile 的构建规则去完成构建. Makefile 由多条规则组成，每条规则都以一个
target（目标）开头，后跟一个 : 冒号，冒号后是这一个目标的 prerequisites（前置条件） 紧接着新的一行，必须以一个 tab 作为开头，后面跟随 command（命令），也就是你希望这一个 target 所执行的构建命令

- target：一个目标代表一条规则，可以是一个或多个文件名。也可以是某个操作的名字（标签），称为伪目标（phony）
- prerequisites：前置条件，这一项是可选参数。通常是多个文件名、伪目标。它的作用是 target 是否需要重新构建的标准，如果前置条件不存在或有过更新（文件的最后一次修改时间）则认为 target 需要重新构建
- command：构建这一个 target 的具体命令集

````makefile
.PHONY: build clean tool lint help

all: build

build:
	go build -v .

tool:
	go tool vet . |& grep -v vendor; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf go-gin-example
	go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
````

1. 在上述文件中，使用了 .PHONY，其作用是声明 build / clean / tool / lint / help 为伪目标，声明为伪目标会怎么样呢？ 声明为伪目标后：在执行对应的命令时，make 就不会去检查是否存在 build
   / clean / tool / lint / help 其对应的文件，而是每次都会运行标签对应的命令 若不声明：恰好存在对应的文件，则 make 将会认为 xx 文件已存在，没有重新构建的必要了

2. 这块比较简单，在命令行执行即可看见效果，实现了以下功能：

- make: make 就是 make all
- make build: 编译当前项目的包和依赖项
- make tool: 运行指定的 Go 工具集
- make lint: golint 一下
- make clean: 删除对象文件和缓存文件
- make help: help

## Question

该项目是完全copy他人的demo用于学习的，代码本身存在一些问题：

- 如命名风格问题，存在大量的字母命名变量或者首字母短语，代码可读性差。
- 耦合性高，作者的架构搭建的风格明显是分层架构模式,但不仅可以看到代码中存在许多上帝类，各层之间互相依赖，没有做到上层只依赖下层原则，
  且框架代码侵入到各个分层中去，如将gin侵入到了util包中，此时明显只需要依赖一个string作为参数即可