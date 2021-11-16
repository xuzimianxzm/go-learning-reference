## Introduction to Web Development

因为Go的net/http包提供了基础的路由函数组合与丰富的功能函数。所以在社区里流行一种用Go编写API不需要框架的观点，在我们看来，如果你的项目的路由在个
位数、URI固定且不通过URI来传递参数，那么确实使用官方库也就足够。但在复杂场景下，官方的http库还是有些力有不逮。

Go的Web框架大致可以分为这么两类:

1. Router框架
2. MVC类框架

在框架的选择上，大多数情况下都是依照个人的喜好和公司的技术栈。例如公司有很多技术人员是PHP出身，那么他们一定会非常喜欢像beego这样的框架，但如果公司有很多C程序员，那么他们的想法可能是越简单越好。比如很多大厂的C程序员甚至可能都会去用C语言去写很小的CGI程序，他们可能本身并没有什么意愿去学习MVC
或者更复杂的Web框架，他们需要的只是一个非常简单的路由（甚至连路由都不需要，只需要一个基础的HTTP协议处理库来帮他省掉没什么意思的体力劳动）。

## Router

Go语言圈子里router也时常被称为http的multiplexer。

### HttpRouter

较流行的开源go Web框架大多使用httprouter，或是基于httprouter的变种对路由进行支持。

因为httprouter中使用的是显式匹配，所以在设计路由的时候需要规避一些会导致路由冲突的情况，例如:

```shell
# conflict:
GET /user/info/:name
GET /user/:id

# no conflict:
GET /user/info/:name
POST /user/:id
```

如果两个路由拥有一致的http方法(指 GET/POST/PUT/DELETE)和请求路径前缀，且在某个位置出现了A路由是wildcard（指:id这种形式）参数，B路由则是普
通字符串，那么就会发生路由冲突。路由冲突会在初始化阶段直接panic:

```shell
panic: wildcard route ':id' conflicts with existing children in path '/user/:id'

goroutine 1 [running]:
...
```

> Notes: 还有一点需要注意，因为httprouter考虑到字典树的深度，在初始化时会对参数的数量进行限制，所以在路由中的参数数目不能超过255，否则会导致httprouter无法识别后续的参数。
> 除支持路径中的wildcard参数之外，httprouter还可以支持*号来进行通配，不过*号开头的参数只能放在路由的结尾，这种设计在RESTful中可能不太常见，主要是为了能够使用httprouter来做简单的HTTP静态文件服务器。

除了正常情况下的路由支持，httprouter也支持对一些特殊情况下的回调函数进行定制，例如404的时候：

```go
r := httprouter.New()
r.NotFound = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
w.Write([]byte("oh no, not found"))
})
```

或者内部panic的时候：

```go
r.PanicHandler = func (w http.ResponseWriter, r *http.Request, c interface{}) {
log.Printf("Recovering from panic, Reason: %#v", c.(error))
w.WriteHeader(http.StatusInternalServerError)
w.Write([]byte(c.(error).Error()))
}
```

目前开源界最为流行（star数最多）的Web框架gin使用的就是httprouter的变种。

### theory

httprouter和众多衍生router使用的数据结构被称为压缩字典树（Radix Tree）。读者可能没有接触过压缩字典树，但对字典树（Trie Tree）应该有所耳闻。下图是一个典型的字典树结构:

![avatar](https://gitee.com/xuzimian/Image/raw/master/Basic/Dictionary_tree.png)

字典树常用来进行字符串检索，例如用给定的字符串序列建立字典树。对于目标字符串，只要从根节点开始深度优先搜索，即可判断出该字符串是否曾经出现过，时间复杂度为O(n)
，n可以认为是目标字符串的长度。为什么要这样做？字符串本身不像数值类型可以进行数值比较，两个字符串对比的时间复杂度取决于字符串长度。如果不用字典树来完成上述功能，要对历史字符串进行排序，再利用二分查找之类的算法去搜索，时间复杂度只高不低。可认为字典树是一种空间换时间的典型做法。

普通的字典树有一个比较明显的缺点，就是每个字母都需要建立一个孩子节点，这样会导致字典树的层数比较深，压缩字典树相对好地平衡了字典树的优点和缺点。是典型的压缩字典树结构

![avatar](https://gitee.com/xuzimian/Image/raw/master/Basic/Compressed_dictionary_tree.png)

每个节点上不只存储一个字母了，这也是压缩字典树中“压缩”的主要含义。使用压缩字典树可以减少树的层数，同时因为每个节点上数据存储也比通常的字典树要多，所以程序的局部性较好（一个节点的path加载到cache即可进行多个字符的对比），从而对CPU缓存友好。