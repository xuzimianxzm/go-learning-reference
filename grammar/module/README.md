## Run Code Command to Create module

> The Dependent module

```shell
cd ./greetings
go mod init xuzimian.com/greetings
```

> Used the Dependent module

```shell
cd  ./hello
go mod init hello
#  after import greetings 'Helo' and use it in hello.go, then run:
go mod edit -replace=xuzimian.com/greetings=../greetings
# run the go mod tidy command to synchronize the xuzimian.com/hello module's dependencies, adding those required by the code,
# but not yet tracked in the module. 该命令主要的作用是整理现有的依赖
go mod tidy

go run hello.go
# or
go build
./hello
# or
go run .
```

## go.mod file

go.mod 文件是启用了 Go modules 的项目所必须的最重要的文件，因为它描述了当前项目（也就是当前模块）的元信息，每一行都以一个动词开头，目前有以下 5 个动词:

- module：用于定义当前项目的模块路径。
- go：用于设置预期的 Go 版本。
- require：用于设置一个特定的模块版本。
- exclude：用于从使用中排除一个特定的模块版本。
- replace：用于将一个模块版本替换为另外一个模块版本。
