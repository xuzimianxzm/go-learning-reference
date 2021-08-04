## The references type slice

When you create a slice, go internally is creating two separate data structures for you. The one is the array data. The
another is what*__* we refer to as the slice, the slices of data structure that has three elements inside of it.

It has a pointer, a capacity number and a length number:

1. The pointer is a pointer over to the underlying array that represents the actual list of items.
2. The capacity is how many elements it can contain at present.
3. The length is represents how many elements currently exist inside the slice.

| Value Types | Reference Types |
|  --------   | --------------- |
|     int     |     slices      |
|    float    |     maps        |
|    string   |     channels    |
|     bool    |     pointers    |
|    structs  |     functions   |

## Array String and slice

Go语言中数组，字符串和切片三者是密切相关的数据结构。这三种数据类型，在底层原始数据有着相同的内存结构，在上层因为语法的限制有着不同的行为表现。

### Array

- 数组是一个由固定长度的特定类型元素组成的序列，且数组长度是数组类型的一部分，不同长度或不同类型的数据组成的数组都是不同的类型，所以无法直接赋值。
- 长度为0的数组(空数组)在内存中并不占用空间，空数组虽然很少直接使用，但是可以用于强调某种特有类型的操作时避免分配额外的内存空间，例如用于通道的同 步操作：
  ```go
  c1 := make(chan [0]int)
  go func(){
  fmt.Println("c1")
  c2 <- [0]int{}
  }
  <- c1
  ```

### Slice

切片的结构和字符串类似，但是解除了制度限制。切片可以和nil比较，只有当切片底层数组为空时切片本身才为nil,这时候切片的len和cap的信息是无效的，如 果有切片底层数据指针为空但len和cap都不为0的情况，那么说明切片本身已经损坏了(
reflect.SliceHeader或unsafe包对切片做了不正确的修改)。

#### 空数组

一般很少用到，但对于切片来说，len为0但cap容量不为0的切片是非常有用的特性。例如下面的TrimSpace()函数用于删除[]byte中的空格，函数利用了 长度为0的切片特性，实现高效且简洁:

````go
func TrimSpace(s []byte) []byte{
b := s[:0]
for _, x := range s {
if x != ' ' {
b = append(b, x)
}
}
return b
}
````

#### Avoid Memory Leak

切片操作并不会复制底层的数据，底层的数组会被保存在内存中，直到它不再被引用。但是有时候可能会因为一个小的内存引用而导致底层整个数组处于被使用的状态， 这会延迟垃圾回收器对底层数组的回收。

例如，FindPhoneNumber()函数加载整个文件到内存，然后搜索第一个出现的电话好吗，最后结果以切片方式返回:

```go
func FindPhoneNumber(filename string) []byte {
b, _ := ioutil.ReadFile(filename)
return regexp.MustCompile("[0-9]+").Find(b)
}
```

这段代码返回的byte[]指向保存整个文件的数组，由于切片引用了整个原始数组，导致垃圾回收器不能及时释放底层数组的空间。要解决这个问题，可以将需要的数据 复制到一个新的切片中(
数据的传值是Go语言编程的一个哲学，虽然传值有一定的代价，但是换取的好处是切断了对原始数据的依赖)。

```go
func FindPhoneNumber(filename string) []byte {
b, _ := ioutil.ReadFile(filename)
b = regexp.MustCompile("[0-9]+").Find(b)
return append([]byte{}, b...)
}
```

#### Slice type forced conversion

为了安全，当两个切片类型[]T 和 []Y 的底层原始切片类型不同时，Go语言时无法直接转换类型的。不过安全都是有一定代价的，有时候这种转换是有它的价值的， 可以转换编码或是提升性能。例如在64系统上，需要对一个[]float64
切片进行高速排序，我们可以将它转换为[]int整数切片，然后以证实的方式进行排序(因为 float64遵循IEEE 754浮点数标准特性，所以当浮点数有序时对应的整数也必然是有序的)。

```go
// +build amd64 arm64

import "sort"

var a = []float64{4, 2, 5, 7, 2, 1, 88, 1}

func SortFloat64FastV1(a []float64) {
// 强制类型转换
var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]

// 以int方式给float64排序
sort.Ints(b)
}

func SortFloat64FastV2(a []float64) {
// 通过 reflect.SliceHeader 更新切片头部信息实现转换
var c []int
aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
*cHdr = *aHdr

// 以int方式给float64排序
sort.Ints(c)
}
```

第一种强制转换是先将切片数据的开始地址转换为一个较大的数组的指针，然后对数组指针对应的数组重新做切片操作。中间需要unsafe.Pointer来连接两个不同类型的指针传递。需要注意的是，Go语言实现中非0大小数组的长度不得超过2GB，因此需要针对数组元素的类型大小计算数组的最大长度范围（[]
uint8最大2GB，[]uint16最大1GB，以此类推，但是[]struct{}数组的长度可以超过2GB）。

第二种转换操作是分别取到两个不同类型的切片头信息指针，任何类型的切片头部信息底层都是对应reflect.SliceHeader结构，然后通过更新结构体方式来更新切片信息，从而实现a对应的[]float64切片到c对应的[]
int类型切片的转换。

通过基准测试，我们可以发现用sort.Ints对转换后的[]int排序的性能要比用sort.Float64s排序的性能好一点。不过需要注意的是，这个方法可行的前提是要保证[]
float64中没有NaN和Inf等非规范的浮点数（因为浮点数中NaN不可排序，正0和负0相等，但是整数中没有这类情形）

### String

Go语言字符串底层数据对应的也是字节数组，但是字符串的制度属性禁止了在程序中对地秤字节数组的元素的修改。字符串赋值只是复制了数据地址和对应的长度，
而不会导致底层数据的复制。字符串虽然不是切片，但是支持切片操作。不同位置的切片底层访问的是同一块内存数据。

Go语言的源文件都采用UTF8编码。因此，Go源文件中出现的字符串面值常量一般也是UTF8编码的(对于转译字符则没有这个限制)，一般都假设Go字符串对应的是 一个合法的UTF8编码的字符序列，可以用for
range循环直接偏离UTF8解码后的Unicode码点值。

## Function Method And Interface

### Function

在Go语言中，函数是第一类对象，我们可以将函数保持到变量中。函数主要有具名和匿名之分，包级函数一般都是具名函数，具名函数是匿名函数的一种特例。

```go
// 具名函数
func Add(a, b int) int {
return a+b
}

// 匿名函数
var Add = func (a, b int) int {
return a+b
}
```

#### Closure

当匿名函数捕获了外部函数的局部变量，这种函数我们一般叫闭包。闭包对捕获的外部变量并不是传值方式访问，而是以引用的方式访问。

```go
func Inc() (v int) {
defer func (){ v++ } ()
return 42
}
```

闭包的这种引用方式访问外部变量的行为可能会导致一些隐含的问题：

```go
func main() {
for i := 0; i < 3; i++ {
defer func (){ println(i) } ()
}
}
// Output:
// 3
// 3
// 3
```

因为是闭包，在for迭代语句中，每个defer语句延迟执行的函数引用的都是同一个i迭代变量，在循环结束后这个变量的值为3，因此最终输出的都是3。 修复的思路是在每轮迭代中为每个defer函数生成独有的变量。可以用下面两种方式：

```go
func main() {
for i := 0; i < 3; i++ {
i := i // 定义一个循环体内局部变量i
defer func (){ println(i) } ()
}
}

func main() {
for i := 0; i < 3; i++ {
// 通过函数传入i
// defer 语句会马上对调用参数求值
defer func (i int){ println(i) } (i)
}
}
```

Note:

1. 一般来说,在for循环内部执行defer语句并不是一个好的习惯，此处仅为示例，不建议使用。

#### pass by value

Go语言中，如果以切片为参数调用函数时，有时候会给人一种参数采用了传引用的方式的假象：因为在被调用函数内部可以修改传入的切片的元素。其实，任何可
以通过函数参数修改调用参数的情形，都是因为函数参数中显式或隐式传入了指针参数。函数参数传值的规范更准确说是只针对数据结构中固定的部分传值，例如字符
串或切片对应结构体中的指针和字符串长度结构体传值，但是并不包含指针间接指向的内容。将切片类型的参数替换为类似reflect.SliceHeader结构体就很好理解 切片传值的含义了:

```go
func twice(x []int) {
for i := range x {
x[i] *= 2
}
}

type IntSliceHeader struct {
Data []int
Len  int
Cap  int
}

func twice(x IntSliceHeader) {
for i := 0; i < x.Len; i++ {
x.Data[i] *= 2
}
}
```

因为切片中的底层数组部分是通过隐式指针传递(指针本身依然是传值的，但是指针指向的却是同一份的数据)，所以被调用函数是可以通过指针修改掉调用参数切片中
的数据。除了数据之外，切片结构还包含了切片长度和切片容量信息，这2个信息也是传值的。如果被调用函数中修改了Len或Cap信息的话，就无法反映到调用参数的
切片中，这时候我们一般会通过返回修改后的切片来更新之前的切片。这也是为何内置的append必须要返回一个切片的原因。

#### Recursion

Go语言函数的递归调用深度逻辑上没有限制，函数调用的栈是不会出现溢出错误的，因为Go语言运行时会根据需要动态地调整函数栈的大小。每个goroutine刚启动
时只会分配很小的栈（4或8KB，具体依赖实现），根据需要动态调整栈的大小，栈最大可以达到GB级（依赖具体实现，在目前的实现中，32位体系结构为250MB,64 位体系结构为1GB）。

### Method

方法一般是面向对象编程(OOP)的一个特性，在C++语言中方法对应一个类对象的成员函数，是关联到具体对象上的虚表中的。但是Go语言的方法却是关联到类型的，
这样可以在编译阶段完成方法的静态绑定。一个面向对象的程序会用方法来表达其属性对应的操作，这样使用这个对象的用户就不需要直接去操作对象，而是借助方法 来做这些事情。

方法是由函数演变而来，只是将函数的第一个对象参数移动到了函数名前面了而已。将第一个函数参数移动到函数前面，从代码角度看虽然只是一个小的改动，但是从 编程哲学角度来看，Go语言已经是进入面向对象语言的行列了。

Go语言不支持传统面向对象中的继承特性，而是以自己特有的组合方式支持了方法的继承。Go语言中，通过在结构体内置匿名的成员来实现继承：

```go
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
Point
Color color.RGBA
}

var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X) // "1"
cp.Point.Y = 2
fmt.Println(cp.Y) // "2"
```

通过嵌入匿名的成员，我们不仅可以继承匿名成员的内部成员，而且可以继承匿名成员类型所对应的方法。我们一般会将Point看作基类，把ColoredPoint看作是它
的继承类或子类。不过这种方式继承的方法并不能实现C++中虚函数的多态特性。所有继承来的方法的接收者参数依然是那个匿名成员本身，而不是当前的变量。

在传统的面向对象语言(eg.C++或Java)的继承中，子类的方法是在运行时动态绑定到对象的，因此基类实现的某些方法看到的this可能不是基类类型对应的对象，
这个特性会导致基类方法运行的不确定性。而在Go语言通过嵌入匿名的成员来“继承”的基类方法，this就是实现该方法的类型的对象，Go语言中方法是编译时静态绑 定的。如果需要虚函数的多态特性，我们需要借助Go语言接口来实现。

### Interface

Go语言之父Rob Pike曾说过一句名言：那些试图避免白痴行为的语言最终自己变成了白痴语言（Languages that try to disallow idiocy become themselves
idiotic）。一般静态编程语言都有着严格的类型系统，这使得编译器可以深入检查程序员有没有作出什么出格的举动。但是，过于严格的类型系统却
会使得编程太过繁琐，让程序员把大好的青春都浪费在了和编译器的斗争中。Go语言试图让程序员能在安全和灵活的编程之间取得一个平衡。它在提供严格的类型检查 的同时，通过接口类型实现了对鸭子类型的支持，使得安全动态的编程变得相对容易。

Go语言中接口类型的独特之处在于它是满足隐式实现的鸭子类型。

所谓鸭子类型说的是：只要走起路来像鸭子、叫起来也像鸭子，那么就可以把它当作鸭子。Go语言中的面向对象就是如此，如果一个对象只要看起来像是某种接口类型
的实现，那么它就可以作为该接口类型使用。这种设计可以让你创建一个新的接口类型满足已经存在的具体类型却不用去破坏这些类型原有的定义；当我们使用的类型
来自于不受我们控制的包时这种设计尤其灵活有用。Go语言的接口类型是延迟绑定，可以实现类似虚函数的多态功能。

#### Implicit conversion

Go语言中，对于基础类型（非接口类型）不支持隐式的转换，我们无法将一个int类型的值直接赋值给int64类型的变量，也无法将int类型的值赋值给底层是int类型的新定义命名类型的变量。Go语言对基础类型的类型一致性要求可谓是非常的严格，但是Go语言对于接口类型的转换则非常的灵活。对象和接口之间的转换、接口和接口之间的转换都可能是隐式的转换。可以看下面的例子：

```go
var (
a io.ReadCloser = (*os.File)(f) // 隐式转换, *os.File 满足 io.ReadCloser 接口
b io.Reader     = a // 隐式转换, io.ReadCloser 满足 io.Reader 接口
c io.Closer = a     // 隐式转换, io.ReadCloser 满足 io.Closer 接口
d io.Reader = c.(io.Reader) // 显式转换, io.Closer 不满足 io.Reader 接口
)
```

#### Private Restriction

有时候对象和接口之间太灵活了，导致我们需要人为地限制这种无意之间的适配。常见的做法是定义一个含特殊方法来区分接口。比如runtime包中的Error接口就定 义了一个特有的RuntimeError方法，用于避免其它类型无意中适配了该接口：

```go
type runtime.Error interface {
error

// RuntimeError is a no-op function but
// serves to distinguish types that are run time
// errors from ordinary errors: a type is a
// run time error if it has a RuntimeError method.
RuntimeError()
}
```

再严格一点的做法是给接口定义一个私有方法。只有满足了这个私有方法的对象才可能满足这个接口，而私有方法的名字是包含包的绝对路径名的，因此只能在包内部 实现这个私有方法才能满足这个接口。测试包中的testing.TB接口就是采用类似的技术：

```go
type testing.TB interface {
Error(args ...interface{})
Errorf(format string, args ...interface{})
...

// A private method to prevent users implementing the
// interface and so future additions to it will not
// violate Go 1 compatibility.
private()
}
```

不过这种通过私有方法禁止外部对象实现接口的做法也是有代价的：首先是这个接口只能包内部使用，外部包正常情况下是无法直接创建满足该接口对象的；其次，这 种防护措施也不是绝对的，恶意的用户依然可以绕过这种保护机制。

在前面的方法一节中我们讲到，通过在结构体中嵌入匿名类型成员，可以继承匿名类型的方法。其实这个被嵌入的匿名成员不一定是普通类型，也可以是接口类型。我
们可以通过嵌入匿名的testing.TB接口来伪造私有的private方法，因为接口方法是延迟绑定，编译时private方法是否真的存在并不重要。

```go
package main

import (
  "fmt"
  "testing"
)

type TB struct {
  testing.TB
}

func (p *TB) Fatal(args ...interface{}) {
  fmt.Println("TB.Fatal disabled!")
}

func main() {
  var tb testing.TB = new(TB)
  tb.Fatal("Hello, playground")
}
```

我们在自己的TB结构体类型中重新实现了Fatal方法，然后通过将对象隐式转换为testing.TB接口类型（因为内嵌了匿名的testing.TB对象，因此是满足
testing.TB接口的），然后通过testing.TB接口来调用我们自己的Fatal方法。

这种通过嵌入匿名接口或嵌入匿名指针对象来实现继承的做法其实是一种纯虚继承，我们继承的只是接口指定的规范，真正的实现在运行的时候才被注入

## Channel

Channel是Go中的一个核心类型，你可以把它看成一个管道，通过它并发核心单元就可以发送或者接收数据进行通讯(communication)。

它的操作符是箭头:<-

```go
ch <- v   // 发送值v到Channel ch中
v := <-ch // 从Channel ch中接收数据，并将数据赋值给v
```

### Channel Tyep

Channel类型的定义格式如下:

```go
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType 
```

它包括三种类型的定义。可选的<-代表channel的方向。如果没有指定方向，那么Channel就是双向的，既可以接收数据，也可以发送数据。

```go
chan T         // 可以接收和发送类型为 T 的数据
chan<- float64 // 只可以用来发送 float64 类型的数据
<-chan int     // 只可以用来接收 int 类型的数据
```

<-总是优先和最左边的类型结合。

```go
chan<- chan int   // 等价 chan<- (chan int)
chan<- <-chan int // 等价 chan<- (<-chan int)
<-chan<-chan int // 等价 <-chan (<-chan int)
chan (<-chan int)
```

容量(capacity)代表Channel容纳的最多的元素的数量，代表Channel的缓存的大小。 如果没有设置容量，或者容量设置为0, 说明Channel没有缓存，只有sender和receiver都准备好了后它们的通讯(
communication)才会发生(Blocking)。 如果设置了缓存，就有可能不发生阻塞， 只有buffer满了后 send才会阻塞， 而只有缓存空了后receive才会阻塞。一个nil channel不会通信。

- 可以在多个goroutine从/往 一个channel 中 receive/send 数据, 不必考虑额外的同步措施。
- Channel可以作为一个先入先出(FIFO)的队列，接收的数据和发送的数据的顺序是一致的。

channel的 receive支持 multi-valued assignment，如:

```go
v, ok := <-ch
```

它可以用来检查Channel是否已经被关闭了,如果OK 是false，表明接收的v是产生的零值，这个channel被关闭了或者为空。

### Send Statement

send语句用来往Channel中发送数据， 如ch <- 3。 它的定义如下:

```go
SendStatement = Channel "<-" Expression
Channel = Expression 
```

在通讯(communication)开始前channel和expression必先求值出来(evaluated)，比如下面的(3+4)先计算出7然后再发送给channel:

```go
c := make(chan int)
defer close(c)
go func () { c <- 3 + 4 }()
i := <-c
fmt.Println(i)
```

send被执行前(proceed)通讯(communication)一直被阻塞着。如前所言，无缓存的channel只有在receiver准备好后send才被执行。如果有缓存，并且缓存 未满，则send会被执行。

- 往一个已经被close的channel中继续发送数据会导致run-time panic。
- 往nil channel中发送数据会一致被阻塞着。

### Receive

<- ch 用来从channel ch中接收数据，这个表达式会一直被block,直到有数据可以接收。

- 从一个nil channel中接收数据会一直被block。
- 从一个被close的channel中接收数据不会被阻塞，而是立即返回，接收完已发送的数据后会返回元素类型的零值(zero value)。

### Blocking

默认情况下，发送和接收会一直阻塞着，直到另一方准备好。这种方式可以用来在gororutine中进行同步，而不必使用显示的锁或者条件变量。

```go
import "fmt"
func sum(s []int, c chan int) {
sum := 0
for _, v := range s {
sum += v
}
c <- sum // send sum to c
}
func main() {
s := []int{7, 2, 8, -9, 4, 0}
c := make(chan int)
go sum(s[:len(s)/2], c)
go sum(s[len(s)/2:], c)
x, y := <-c, <-c // 这句会一直等待计算结果发送到channel中。
fmt.Println(x, y, x+y)
}
```

### Buffered Channels

make的第二个参数指定缓存的大小：ch := make(chan int, 100)。 通过缓存的使用，可以尽量避免阻塞，提供应用的性能。

### Close

内建的close方法可以用来关闭channel。 总结一下channel关闭后sender的receiver操作。

- 如果channel c已经被关闭,继续往它发送数据会导致panic: send on closed channel:

```go
import "time"
func main() {
go func () {
time.Sleep(time.Hour)
}()
c := make(chan int, 10)
c <- 1
c <- 2
close(c)
c <- 3
}
```

但是从这个关闭的channel中不但可以读取出已发送的数据，还可以不断的读取零值:

```go
c := make(chan int, 10)
c <- 1
c <- 2
close(c)
fmt.Println(<-c) //1
fmt.Println(<-c) //2
fmt.Println(<-c) //0
fmt.Println(<-c) //0
```

### Range

for range语句可以处理Channel:

```go
func main() {
go func () {
time.Sleep(1 * time.Hour)
}()
c := make(chan int)
go func () {
for i := 0; i < 10; i = i + 1 {
c <- i
}
close(c)
}()
for i := range c {
fmt.Println(i)
}
fmt.Println("Finished")
}
```

range c产生的迭代值为Channel中发送的值，它会一直迭代直到channel被关闭。上面的例子中如果把close(c)注释掉，程序会一直阻塞在for range那一行。

### Select

select语句选择一组可能的send操作和receive操作去处理。它类似switch,但是只是用来处理通讯(communication)操作。 它的case可以是send语句，也可以是receive语句，亦或者default。

```go
import "fmt"
func fibonacci(c, quit chan int) {
x, y := 0, 1
for {
select {
case c <- x:
x, y = y, x+y
case <-quit:
fmt.Println("quit")
return
}
}
}
func main() {
c := make(chan int)
quit := make(chan int)
go func () {
for i := 0; i < 10; i++ {
fmt.Println(<-c)
}
quit <- 0
}()
fibonacci(c, quit)
}
```

如果有同时多个case去处理,比如同时有多个channel可以接收数据，那么Go会伪随机的选择一个case处理(pseudo-random)。如果没有case需要处理，则会选 择default去处理，如果default
case存在的情况下。如果没有default case，则select语句会阻塞，直到某个case需要处理。

需要注意的是，nil channel上的操作会一直被阻塞，如果没有default case,只有nil channel的select会一直被阻塞。

- select语句和switch语句一样，它不是循环，它只会选择一个case来处理，如果想一直处理channel，可以在外面加一个无限的for循环：

```go
for {
select {
case c <- x:
x, y = y, x+y
case <-quit:
fmt.Println("quit")
return
}
}
```

### Timeout

select有很重要的一个应用就是超时处理。 因为上面我们提到，如果没有case需要处理，select语句就会一直阻塞着。这时候我们可能就需要一个超时操作，用 来处理超时的情况。

下面这个例子我们会在2秒后往channel c1中发送一个数据，但是select设置为1秒超时,因此我们会打印出timeout 1,而不是result 1。

```go
import "time"
import "fmt"
func main() {
c1 := make(chan string, 1)
go func () {
time.Sleep(time.Second * 2)
c1 <- "result 1"
}()
select {
case res := <-c1:
fmt.Println(res)
case <-time.After(time.Second * 1):
fmt.Println("timeout 1")
}
}
```

其实它利用的是time.After方法，它返回一个类型为<-chan Time的单向的channel，在指定的时间发送一个当前时间给返回的channel中。

### Timer And Ticker

timer是一个定时器，代表未来的一个单一事件，你可以告诉timer你要等待多长时间，它提供一个Channel，在将来的那个时间那个Channel提供了一个时间值。 下面的例子中第二行会阻塞2秒钟左右的时间，直到时间到了才会继续执行。

```go
timer1 := time.NewTimer(time.Second * 2)
<-timer1.C
fmt.Println("Timer 1 expired")
```

当然如果你只是想单纯的等待的话，可以使用time.Sleep来实现。

还可以使用timer.Stop来停止计时器。

```go
timer2 := time.NewTimer(time.Second)
go func () {
<-timer2.C
fmt.Println("Timer 2 expired")
}()
stop2 := timer2.Stop()
if stop2 {
fmt.Println("Timer 2 stopped")
}
```

ticker是一个定时触发的计时器，它会以一个间隔(interval)往Channel发送一个事件(当前时间)，而Channel的接收者可以以固定的时间间隔从Channel中读
取事件。下面的例子中ticker每500毫秒触发一次，你可以观察输出的时间。

```go
ticker := time.NewTicker(time.Millisecond * 500)
go func () {
for t := range ticker.C {
fmt.Println("Tick at", t)
}
}()
```

## Concurrence-Oriented Programming

常见的并行编程有多种模型，主要有多线程、消息传递等。从理论上来看，多线程和基于消息的并发编程是等价的。由于多线程并发模型可以自然对应到多核的处理器，
主流的操作系统因此也都提供了系统级的多线程支持，同时从概念上讲多线程似乎也更直观，因此多线程编程模型逐步被吸纳到主流的编程语言特性或语言扩展库中。
而主流编程语言对基于消息的并发编程模型支持则相比较少，Erlang语言是支持基于消息传递并发编程模型的代表者，它的并发体之间不共享内存。 Go语言是基于 消息并发模型的集大成者，它将基于
CSP模型的并发编程内置到了语言中，通过一个go关键字就可以轻易地启动一个Goroutine，与Erlang不同的是Go语言的 Goroutine之间是共享内存的。

### Goroutine And System Thread

Goroutine是Go语言特有的并发体，是一种轻量级的线程，由go关键字启动。在真实的Go语言的实现中，goroutine和系统线程也不是等价的。尽管两者的区别实 际上只是一个量的区别，但正是这个量变引发了Go语言并发编程质的飞跃。

首先，每个系统级线程都会有一个固定大小的栈（一般默认可能是2MB），这个栈主要用来保存函数递归调用时参数和局部变量。固定了栈的大小导致了两个问题：一
是对于很多只需要很小的栈空间的线程来说是一个巨大的浪费，二是对于少数需要巨大栈空间的线程来说又面临栈溢出的风险。针对这两个问题的解决方案是：要么降
低固定的栈大小，提升空间的利用率；要么增大栈的大小以允许更深的函数递归调用，但这两者是没法同时兼得的。相反，一个Goroutine会以一个很小的栈启动（
可能是2KB或4KB），当遇到深度递归导致当前栈空间不足时，Goroutine会根据需要动态地伸缩栈的大小（主流实现中栈的最大值可达到1GB）。因为启动的代价很 小，所以我们可以轻易地启动成千上万个Goroutine。

Go的运行时还包含了其自己的调度器，这个调度器使用了一些技术手段，可以在n个操作系统线程上多工调度m个Goroutine。Go调度器的工作和内核的调度是相似
的，但是这个调度器只关注单独的Go程序中的Goroutine。Goroutine采用的是半抢占式的协作调度，只有在当前Goroutine发生阻塞时才会导致调度；同时发生
在用户态，调度器会根据具体函数只保存必要的寄存器，切换的代价要比系统线程低得多。运行时有一个runtime.GOMAXPROCS变量，用于控制当前运行正常非阻塞 Goroutine的系统线程数目。

### Atomic Operation

所谓的原子操作就是并发编程中“最小的且不可并行化”的操作。通常，如果多个并发体对同一个共享资源进行的操作是原子的话，那么同一时刻最多只能有一个并发体
对该资源进行操作。从线程角度看，在当前线程修改共享资源期间，其它的线程是不能访问该资源的。原子操作对于多线程并发编程模型来说，不会发生有别于单线程 的意外情况，共享资源的完整性可以得到保证。

一般情况下，原子操作都是通过“互斥”访问来保证的，通常由特殊的CPU指令提供保护。当然，如果仅仅是想模拟下粗粒度的原子操作，我们可以借助于sync.Mutex 来实现：

```go
import (
"sync"
)

var total struct {
sync.Mutex
value int
}

func worker(wg *sync.WaitGroup) {
defer wg.Done()

for i := 0; i <= 100; i++ {
total.Lock()
total.value += i
total.Unlock()
}
}

func main() {
var wg sync.WaitGroup
wg.Add(2)
go worker(&wg)
go worker(&wg)
wg.Wait()

fmt.Println(total.value)
}
```

用互斥锁来保护一个数值型的共享资源，麻烦且效率低下。标准库的sync/atomic包对原子操作提供了丰富的支持。我们可以重新实现上面的例子：

```go
func worker(wg *sync.WaitGroup) {
defer wg.Done()

var i uint64
for i = 0; i <= 100; i++ {
atomic.AddUint64(&total, i)
}
}
```

### Sequential consistent memory model

```go
var a string
var done bool

func setup() {
a = "hello, world"
done = true
}

func main() {
go setup()
for !done {}
print(a)
}
```

我们创建了setup线程，用于对字符串a的初始化工作，初始化完成之后设置done标志为true。main函数所在的主线程中，通过for !done {}检测done变为true 时，认为字符串初始化工作完成，然后进行字符串的打印工作。

但是Go语言并不保证在main函数中观测到的对done的写入操作发生在对字符串a的写入的操作之后，因此程序很可能打印一个空字符串。更糟糕的是，因为两个线程
之间没有同步事件，setup线程对done的写入操作甚至无法被main线程看到，main函数有可能陷入死循环中。

在Go语言中，同一个Goroutine线程内部，顺序一致性内存模型是得到保证的。但是不同的Goroutine之间，并不满足顺序一致性内存模型，需要通过明确定义的同
步事件来作为同步的参考。如果两个事件不可排序，那么就说这两个事件是并发的。为了最大化并行，Go语言的编译器和处理器在不影响上述规定的前提下可能会对执 行语句重新排序（CPU也会对一些指令进行乱序执行）。

比如下面这个程序：

```go
func main() {
go println("你好, 世界")
}
```

根据Go语言规范，main函数退出时程序结束，不会等待任何后台线程。因为Goroutine的执行和main函数的返回事件是并发的，谁都有可能先发生，所以什么时候 打印，能否打印都是未知的。

用前面的原子操作并不能解决问题，因为我们无法确定两个原子操作之间的顺序。解决问题的办法就是通过同步原语来给两个事件明确排序：

```go
func main() {
done := make(chan int)

go func(){
println("你好, 世界")
done <- 1
}()

<-done
}
```

当<-done执行时，必然要求done <- 1也已经执行。根据同一个Gorouine依然满足顺序一致性规则，我们可以判断当done <- 1执行时， println("你好, 世界")
语句必然已经执行完成了。因此，现在的程序确保可以正常打印结果。

### Initialization Sequence

Go程序的初始化和执行总是从main.main函数开始的。但是如果main包里导入了其它的包，则会按照顺序将它们包含进main包里（这里的导入顺序依赖具体实现，
一般可能是以文件名或包路径名的字符串顺序导入）。如果某个包被多次导入的话，在执行的时候只会导入一次。当一个包被导入时，如果它还导入了其它的包，则先
将其它的包包含进来，然后创建和初始化这个包的常量和变量。然后就是调用包里的init函数，如果一个包有多个init函数的话，实现可能是以文件名的顺序调用，
同一个文件内的多个init则是以出现的顺序依次调用（init不是普通函数，可以定义有多个，所以不能被其它函数调用）。最终，在main包的所有包常量、包变量被
创建和初始化，并且init函数被执行后，才会进入main.main函数，程序开始正常执行。下图是Go程序函数启动顺序的示意图：
![avatar](https://gitee.com/xuzimian/Image/raw/master/golang/go-Initialization-Sequence.png)

- 在main.main函数执行之前所有代码都运行在同一个Goroutine中，也是运行在程序的主系统线程中。如果某个init函数内部用go关键字启动了新的Goroutine 的话，新的Goroutine和main.main函数是并发执行的

### Goroutine Creating

go语句会在当前Goroutine对应函数返回前创建新的Goroutine. 例如:

```go
var a string

func f() {
print(a)
}

func hello() {
a = "hello, world"
go f()
}
```

执行go f()语句创建Goroutine和hello函数是在同一个Goroutine中执行, 根据语句的书写顺序可以确定Goroutine的创建发生在hello函数返回之前, 但 是 新创建Goroutine对应的f()
的执行事件和hello函数返回的事件则是不可排序的，也就是并发的。调用hello可能会在将来的某一时刻打印"hello, world"， 也很可能是在hello函数执行完成后才打印。

### Channel-based Communication

Channel通信是在Goroutine之间进行同步的主要方法。在无缓存的Channel上的每一次发送操作都有与其对应的接收操作相配对，发送和接收操作通常发生在不
同的Goroutine上（在同一个Goroutine上执行2个操作很容易导致死锁）。

- [无缓存的]Channel上的发送操作总在对应的接收操作 [完成前] 发生,对于从[无缓存]Channel进行的接收，发生在对该Channel进行的发送 [完成之前]。

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
msg = "你好, 世界"
done <- true
}

func main() {
go aGoroutine()
<-done
println(msg)
}
```

可保证打印出“hello, world”。该程序首先对msg进行写入，然后在done管道上发送同步信号，随后从done接收对应的同步信号，最后执行println函数。

若在关闭Channel后继续从中接收数据，接收者就会收到该Channel返回的零值。因此在这个例子中，用close(c)关闭管道代替done <- false依然能保证该程 序产生相同的行为。

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
msg = "你好, 世界"
close(done)
}

func main() {
go aGoroutine()
<-done
println(msg)
}
```

对于 [从] 无缓存Channel进行的接收(注意「非操作」)，发生在对该Channel进行的发送 [完成之前]。

基于上面这个规则可知，交换两个Goroutine中的接收和发送操作也是可以的（但是很危险):

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
msg = "hello, world"
<- done // 从channel接收,这时候channel 由于还没有接收数据，所以这个 channel 是无缓存的，它会发生在对这个channel 发送数据之前，所以会block。
}
func main() {
go aGoroutine()
done <- true // 发送消息到 channel
println(msg)
}
```

### Common Concurrence Model

首先要明确一个概念：并发不是并行。并发更关注的是程序的设计层面，并发的程序完全是可以顺序执行的，只有在真正的多核CPU上才可能真正地同时运行。并行更
关注的是程序的运行层面，并行一般是简单的大量重复，例如GPU中对图像处理都会有大量的并行运算。

