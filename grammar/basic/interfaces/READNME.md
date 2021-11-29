## syntax

使用 type 关键字来申明，interface 代表类型，大括号里面定义接口的方法签名集合。

````go
type Animal interface {
Bark() string
Walk() string
}

type Dog struct {}

func (d Dog) Bark() string {}

func (d Dog) Walk() string {}

````

> Note: 必须是同时实现 Bark() 和Walk() 方法，否则都不能算实现了Animal接口。

### function implement interface

Go语言中的所有类型都可以实现接口,函数作为Go语言中的一种类型，也不例外。

```go
type Printer interface {
  Print(p interface{})
}

type FuncCaller func (p interface{})

func (funcCaller FuncCaller) Print(p interface{}) {
   funcCaller(p)
}

func main() {
	var printer Printer
	printer = FuncCaller(func (p interface{}) {
		fmt.Println(p)
    })
	
	printer.Print("Any thing");
}

```

### nil interface

官方定义：Interface values with nil underlying values:

- 只声明没赋值的interface 是nil interface，value和 type 都是 nil
- 只要赋值了，即使赋了一个值为nil类型，也不再是nil interface

### empty interface

Go 允许不带任何方法的 interface ,这种类型的 interface 叫 empty interface。所有类型都实现了 empty interface,因为任何一种类型至少实现了 0 个方法。 典型的应用场景是
fmt包的Print方法，它能支持接收各种不同的类型的数据，并且输出到控制台,就是interface{}的功劳。

```go
func Print(i interface{}) {
fmt.Println(i)
}
```

### determine interface type

一个 interface 可被多种类型实现，有时候我们需要区分 interface 变量究竟存储哪种类型的值？类型断言提供对接口值的基础具体值的访问

```go
t := i.(T)
```

该语句断言接口值i保存的具体类型为T，并将T的基础值分配给变量t。如果i保存的值不是类型 T ，将会触发 panic 错误。为了避免 panic 错误发生，可以通过 如下操作来进行断言检查

```go
t, ok := i.(T)
```

断言成功，ok 的值为 true,断言失败 t 值为T类型的零值,并且不会发生 panic 错误。

### type switches

```go
switch v := i.(type) {
case T:
// here v has type T
case S:
// here v has type S
default:
// no match; here v has the same type as i
}
```

### Stringers

One of the most ubiquitous interfaces is Stringer defined by the fmt package.

```go
type Stringer interface {
String() string
}
```

A Stringer is a type that can describe itself as a string. The fmt package (and many others) look for this interface to
print values.

### Errors

Go 语言通过内置的错误接口提供了非常简单的错误处理机制。 error类型是一个接口类型，这是它的定义：

```go
type error interface {
Error() string
}
```

我们可以在编码中通过实现 error 接口类型来生成错误信息。 函数通常在最后的返回值中返回错误信息。使用errors.New 可返回一个错误信息：

```go
func Sqrt(f float64) (float64, error) {
if f < 0 {
return 0, errors.New("math: square root of negative number")
}
// 实现
}
```