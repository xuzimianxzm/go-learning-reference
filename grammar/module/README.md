### Run Code Command to Create module

```shell
cd ./greetings
go mod init xuzimian.com/greetings
```

```shell
cd  ./hello
go mod init hello
#  after import greetings 'Helo' and use it in hello.go, then run:
go mod edit -replace=xuzimian.com/greetings=../greetings
# run the go mod tidy command to synchronize the xuzimian.com/hello module's dependencies, adding those required by the code,
# but not yet tracked in the module.
go mod tidy

go run hello.go 
# or 
go build
./hello 
# or 
go run .
```