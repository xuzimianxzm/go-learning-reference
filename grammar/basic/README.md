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
