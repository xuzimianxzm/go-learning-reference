## Formatting

The go provides the commands program of gofmt to formatting the go code style. which operates at the package level
rather than source file level.

```shell
go fmt
```

## Commentary

Go provides C-style /\* \*/ block comments and C++-style // line comments.

- Every package should have a package comment,a block comment preceding the package clause.
- For multi-file packages, the package comment only needs to be present in one file, and any one will do.
- If the package is simple, the package comment can be brief.
- Every exported (capitalized) name in a program should have a doc comment.
- The first sentence should be a one-sentence summary that starts with the name being declared.

If every doc comment begins with the name of the item it describes, you can use the doc subcommand of the go tool and
run the output through grep.

```shell
go doc -all regexp | grep -i item-name
```

## Package names

Notes: Don't use the import . notation, which can simplify tests that must run outside the package they are testing, but
should otherwise be avoided.

## Interface names

- one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent
  noun: Reader, Writer, Formatter, CloseNotifier etc.

## Named result parameters

The return or result "parameters" of a Go function can be given names and used as regular variables, just like the
incoming parameters. When named, they are initialized to the zero values for their types when the function begins; if
the function executes a return statement with no arguments, the current values of the result parameters are used as the
returned values.

```go
func ReadFull(r Reader, buf []byte) (n int, err error) {
for len(buf) > 0 && err == nil {
var nr int
nr, err = r.Read(buf)
n += nr
buf = buf[nr:]
}
return
}
```

## Defer

Go's defer statement schedules a function call (the deferred function) to be run immediately before the function
executing the defer returns. It's an unusual but effective way to deal with situations such as resources that must be
released regardless of which path a function takes to return. The canonical examples are unlocking a mutex or closing a
file.

```go
// Contents returns the file's contents as a string.
func Contents(filename string) (string, error) {
f, err := os.Open(filename)
if err != nil {
return "", err
}
defer f.Close() // f.Close will run when we're finished.

var result []byte
buf := make([]byte, 100)
for {
n, err := f.Read(buf[0:])
result = append(result, buf[0:n]...) // append is discussed later.
if err != nil {
if err == io.EOF {
break
}
return "", err // f will be closed if we return here.
}
}
return string(result), nil // f will be closed if we return here.
}
```

Deferring a call to a function such as Close has two advantages. First, it guarantees that you will never forget to
close the file, a mistake that's easy to make if you later edit the function to add a new return path. Second, it means
that the close sits near the open, which is much clearer than placing it at the end of the function.

The arguments to the deferred function (which include the receiver if the function is a method) are evaluated when the
defer executes, not when the call executes. Besides avoiding worries about variables changing values as the function
executes, this means that a single deferred call site can defer multiple function executions. Here's a silly example.

```go
for i := 0; i < 5; i++ {
defer fmt.Printf("%d ", i)
}
```

Deferred functions are executed in LIFO order, so this code will cause 4 3 2 1 0 to be printed when the function
returns.

## Data

### Allocation with new

Go has two allocation primitives, the built-in functions new and make. They do different things and apply to different
types, which can be confusing, but the rules are simple. Let's talk about new first. It's a built-in function that
allocates memory, but unlike its namesakes in some other languages it does not initialize the memory, it only zeros it.
That is, new(T) allocates zeroed storage for a new item of type T and returns its address, a value of type \*T.

- Since the memory returned by new is zeroed, it's helpful to arrange when designing your data structures that the zero
  value of each type can be used without further initialization.

The zero-value-is-useful property works transitively. Consider this type declaration.

```go

type SyncedBuffer struct {
lock    sync.Mutex
buffer  bytes.Buffer
}

/*
   Values of type SyncedBuffer are also ready to use immediately upon allocation or just declaration. In the next snippet, both p and v will work correctly without further arrangement.
*/
p := new(SyncedBuffer) // type *SyncedBuffer
var v SyncedBuffer // type  SyncedBuffer

```

### Constructors and composite literals

The fields of a composite literal are laid out in order and must all be present. However, by labeling the elements
explicitly as field:value pairs, the initializers can appear in any order, with the missing ones left as their
respective zero values. Thus we could say.

```go
File{fd: fd, name: name}
```

- As a limiting case, if a composite literal contains no fields at all, it creates a zero value for the type. The
  expressions new(File) and &File{} are equivalent.

### Allocation with make

The built-in function make(T, args) serves a purpose different from new(T). It creates slices, maps, and channels only,
and it returns an initialized (not zeroed) value of type T (not \*T). The reason for the distinction is that these three
types represent, under the covers, references to data structures that must be initialized before use.

A slice, for example, is a three-item descriptor containing a pointer to the data (inside an array), the length, and the
capacity, and until those items are initialized, the slice is nil.

These examples illustrate the difference between new and make.

```go
var p *[]int = new([]int) // allocates slice structure; *p == nil; rarely useful
var v  []int = make([]int, 100) // the slice v now refers to a new array of 100 ints

// Unnecessarily complex:
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// Idiomatic:
v := make([]int, 100)
```

### Arrays

There are major differences between the ways arrays work in Go and C. In Go,

- Arrays are values. Assigning one array to another copies all the elements.
- In particular, if you pass an array to a function, it will receive a copy of the array, not a pointer to it.
- The size of an array is part of its type. The types [10]int and [20]int are distinct.

### Slices

- We must return the slice afterwards because, although Append can modify the elements of slice, the slice itself (the
  run-time data structure holding the pointer, length, and capacity) is passed by value.

### Two-dimensional slices

- Go's arrays and slices are one-dimensional. To create the equivalent of a 2D array or slice, it is necessary to define
  an array-of-arrays or slice-of-slices, like this:

```go
type Transform [3][3]float64 // A 3x3 array, really an array of arrays.
type LinesOfText [][]byte // A slice of byte slices.
```

### Maps

An attempt to fetch a map value with a key that is not present in the map will return the zero value for the type of the
entries in the map.For instance, if the map contains integers, looking up a non-existent key will return 0.

Sometimes you need to distinguish a missing entry from a zero value. Is there an entry for "UTC" or is that 0 because
it's not in the map at all? You can discriminate with a form of multiple assignment.

```go
var seconds int
var ok bool
seconds, ok = timeZone[tz]
```

For obvious reasons this is called the “comma ok” idiom. In this example, if tz is present, seconds will be set
appropriately and ok will be true; if not, seconds will be set to zero and ok will be false.

## Printing

Our String method is able to call Sprintf because the print routines are fully reentrant and can be wrapped this way.
There is one important detail to understand about this approach, however: don't construct a String method by calling
Sprintf in a way that will recur into your String method indefinitely.

```go
type MyString string

func (m MyString) String() string {
return fmt.Sprintf("MyString=%s", m) // Error: will recur forever.
}
```

It's also easy to fix: convert the argument to the basic string type, which does not have the method.

```go
type MyString string
func (m MyString) String() string {
return fmt.Sprintf("MyString=%s", string(m)) // OK: note conversion.
}
```

## Append

```go
func append(slice []T, elements ...T) []T
```

where T is a placeholder for any given type. You can't actually write a function in Go where the type T is determined by
the caller. That's why append is built in: it needs support from the compiler.

```go
// append element
x := []int{1, 2, 3}
x = append(x, 4, 5, 6)
fmt.Println(x)

// or append slice
x := []int{1, 2, 3}
y := []int{4, 5, 6}
x = append(x, y...)
fmt.Println(x)
```

## Constants

They are created at compile time, even when defined as locals in functions, and can only be numbers, characters (runes),
strings or booleans. Because of the compile-time restriction, the expressions that define them must be constant
expressions, evaluatable by the compiler. For instance, 1<<3 is a constant expression, while math.Sin(math.Pi/4) is not
because the function call to math.Sin needs to happen at run time.

In Go, enumerated constants are created using the iota enumerator. Since iota can be part of an expression and
expressions can be implicitly repeated, it is easy to build intricate sets of values.

```go
type ByteSize float64

const (
_ = iota // ignore first value by assigning to blank identifier
KB ByteSize = 1 << (10 * iota)
MB
GB
TB
PB
EB
ZB
YB
);

func (b ByteSize) String() string {
switch {
case b >= YB:
return fmt.Sprintf("%.2fYB", b/YB)
case b >= ZB:
return fmt.Sprintf("%.2fZB", b/ZB)
case b >= EB:
return fmt.Sprintf("%.2fEB", b/EB)
case b >= PB:
return fmt.Sprintf("%.2fPB", b/PB)
case b >= TB:
return fmt.Sprintf("%.2fTB", b/TB)
case b >= GB:
return fmt.Sprintf("%.2fGB", b/GB)
case b >= MB:
return fmt.Sprintf("%.2fMB", b/MB)
case b >= KB:
return fmt.Sprintf("%.2fKB", b/KB)
}
return fmt.Sprintf("%.2fB", b)
}
```

## 类型命名和类型声明的区别

- 类型别名的语法: type identifier = Type
- 类型定义的语法: type type-name type-underlying
- 类型别名和原类型是相同的，而类型定义和原类型是不同的两个类型。
  > 完全一样(identical types)意味着这两种类型的数据可以互相赋值，而类型定义要和原始类型赋值的时候需要类型转换(Conversion T(x))。

### 类型循环

类型别名在定义的时候不允许出现循环定义别名的情况，如下面所示：

```go
type T1 = T2
type T2 = T1

// or:
type T1 = struct {
next *T2
}
type T2 = T1
```

### 方法集

既然类型别名和原始类型是相同的，那么它们的方法集也是相同的,下面的例子中 T1 和 T3 都有 say 和 greeting 方法:

```go
type T1 struct{}
type T3 = T1
func (t1 T1) say(){}
func (t3 *T3) greeting(){}
func main() {
var t1 T1
// var t2 T2
var t3 T3
t1.say()
t1.greeting()
t3.say()
t3.greeting()
}
```

如果类型别名和原始类型定义了相同的方法，代码编译的时候会报错，因为有重复的方法定义。

### byte 和 rune 类型

在 Go 1.9 中， 内部其实使用了类型别名的特性。 比如内建的 byte 类型，其实是 uint8 的类型别名，而 rune 其实是 int32 的类型别名。

## The init function

Finally, each source file can define its own niladic init function to set up whatever state is required. (Actually each
file can have multiple init functions.) And finally means finally: init is called after all the variable declarations in
the package have evaluated their initializers, and those are evaluated only after all the imported packages have been
initialized.

Besides initializations that cannot be expressed as declarations, a common use of init functions is to verify or repair
correctness of the program state before real execution begins.

```go
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
    // gopath may be overridden by --gopath flag on command line.
    flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```

## Methods

### Pointers vs. Values

- The rule about pointers vs. values for receivers is that value methods can be invoked on pointers and values, but pointer methods can only be invoked on pointers.This rule arises because pointer methods can modify the receiver; invoking them on a value would cause the method to receive a copy of the value, so any modifications would be discarded.The language therefore disallows this mistake.

### Interfaces and other types

#### Interfaces

Interfaces in Go provide a way to specify the behavior of an object: if something can do this, then it can be used here.Interfaces with only one or two methods are common in Go code, and are usually given a name derived from the method, such as io.Writer for something that implements Write.

```go
type Sequence []int

// Methods required by sort.Interface.
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// Copy returns a copy of the Sequence.
func (s Sequence) Copy() Sequence {
    copy := make(Sequence, 0, len(s))
    return append(copy, s...)
}

// Method for printing - sorts the elements before printing.
func (s Sequence) String() string {
    s = s.Copy() // Make a copy; don't overwrite argument.
    sort.Sort(s)
    str := "["
    for i, elem := range s { // Loop is O(N²); will fix that in next example.
        if i > 0 {
            str += " "
        }
        str += fmt.Sprint(elem)
    }
    return str + "]"
}
```

#### Conversions

The String method of Sequence(The above example) is recreating the work that Sprint already does for slices. (It also has complexity O(N²), which is poor.) We can share the effort (and also speed it up) if we convert the Sequence to a plain []int before calling Sprint.

```go
func (s Sequence) String() string {
    s = s.Copy()
    sort.Sort(s)
    return fmt.Sprint([]int(s))
}
```

This method is another example of the conversion technique for calling Sprintf safely from a String method. Because the two types (Sequence and []int) are the same if we ignore the type name, it's legal to convert between them. The conversion doesn't create a new value, it just temporarily acts as though the existing value has a new type. (There are other legal conversions, such as from integer to floating point, that do create a new value.)

```go
type Sequence []int

// Method for printing - sorts the elements before printing
func (s Sequence) String() string {
    s = s.Copy()
    sort.IntSlice(s).Sort()
    return fmt.Sprint([]int(s))
}
```

Now, instead of having Sequence implement multiple interfaces (sorting and printing), we're using the ability of a data item to be converted to multiple types (Sequence, sort.IntSlice and []int), each of which does some part of the job. That's more unusual in practice but can be effective.

#### Interface conversions and type assertions

Type switches are a form of conversion: they take an interface and, for each case in the switch, in a sense convert it to the type of that case. Here's a simplified version of how the code under fmt.Printf turns a value into a string using a type switch. If it's already a string, we want the actual string value held by the interface, while if it has a String method we want the result of calling the method.

```go
type Stringer interface {
    String() string
}

var value interface{} // Value provided by caller.
switch str := value.(type) {
case string:
    return str
case Stringer:
    return str.String()
}
```

The first case finds a concrete value; the second converts the interface into another interface. It's perfectly fine to mix types this way.

What if there's only one type we care about? If we know the value holds a string and we just want to extract it? A one-case type switch would do, but so would a type assertion. A type assertion takes an interface value and extracts from it a value of the specified explicit type. The syntax borrows from the clause opening a type switch, but with an explicit type rather than the type keyword:

```go
value.(typeName)
```

and the result is a new value with the static type typeName. That type must either be the concrete type held by the interface, or a second interface type that the value can be converted to. To extract the string we know is in the value, we could write:

```go
str := value.(string)
```

As an illustration of the capability, here's an if-else statement that's equivalent to the type switch that opened this section.

```go
if str, ok := value.(string); ok {
    return str
} else if str, ok := value.(Stringer); ok {
    return str.String()
}
```

#### Generality

If a type exists only to implement an interface and will never have exported methods beyond that interface, there is no need to export the type itself. Exporting just the interface makes it clear the value has no interesting behavior beyond what is described in the interface. It also avoids the need to repeat the documentation on every instance of a common method.

```go
/* Ver1 与 Ver2 是 Version 的两种不同的具体实现。
 * NewShow 接受一个 Version 作为参数，并返回一个 Show。
 * NewShow 不管 Version 的具体实现，只要是 Version 都可以接受，
 *         并且不管 Show 的具体实现，本例子虽然只给出了一种 Show 的实现 (ShowVer)
 *         但事实上换另一种实现也是可以的。
 */

package main

import "fmt"

type Version interface {
    Is() string
}

type Show interface {
    VerNum()
}

type ShowVer struct{
    V string
}
func (s ShowVer) VerNum() {
    fmt.Println("Ver.", s.V)
}

type Ver1 struct {
    V string
}
func (v Ver1) Is() string {
    return fmt.Sprint(v.V)
}

type Ver2 struct{}
func (v Ver2) Is() string {
    return "2"
}

func main() {
    ver1 := Ver1{"1"}
    ver2 := Ver2{}
    show1 := NewShow(ver1)
    show2 := NewShow(ver2)
    show1.VerNum()
    show2.VerNum()
}

func NewShow(ver Version) Show {
    v := ver.Is()
    return ShowVer{v}
}
```

#### Interfaces and methods

Since almost anything can have methods attached, almost anything can satisfy an interface. One illustrative example is in the http package, which defines the Handler interface. Any object that implements Handler can serve HTTP requests.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

ResponseWriter is itself an interface that provides access to the methods needed to return the response to the client. Those methods include the standard Write method, so an http.ResponseWriter can be used wherever an io.Writer can be used. Request is a struct containing a parsed representation of the request from the client.

For brevity, let's ignore POSTs and assume HTTP requests are always GETs; that simplification does not affect the way the handlers are set up. Here's a trivial implementation of a handler to count the number of times the page is visited.

```go
// Simple counter server.
type Counter struct {
    n int
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ctr.n++
    fmt.Fprintf(w, "counter = %d\n", ctr.n)
}
```

(Keeping with our theme, note how Fprintf can print to an http.ResponseWriter.) In a real server, access to ctr.n would need protection from concurrent access. See the sync and atomic packages for suggestions.

For reference, here's how to attach such a server to a node on the URL tree.

```go
import "net/http"
...
ctr := new(Counter)
http.Handle("/counter", ctr)
```

But why make Counter a struct? An integer is all that's needed. (The receiver needs to be a pointer so the increment is visible to the caller.)

```go
// Simpler counter server.
type Counter int

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    *ctr++
    fmt.Fprintf(w, "counter = %d\n", *ctr)
}
```

What if your program has some internal state that needs to be notified that a page has been visited? Tie a channel to the web page.

```go
// A channel that sends a notification on each visit.
// (Probably want the channel to be buffered.)
type Chan chan *http.Request

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ch <- req
    fmt.Fprint(w, "notification sent")
}
```

Finally, let's say we wanted to present on /args the arguments used when invoking the server binary. It's easy to write a function to print the arguments.

```go
func ArgServer() {
    fmt.Println(os.Args)
}
```

How do we turn that into an HTTP server? We could make ArgServer a method of some type whose value we ignore, but there's a cleaner way. Since we can define a method for any type except pointers and interfaces, we can write a method for a function. The http package contains this code:

```go
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, req).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
    f(w, req)
}
```

HandlerFunc is a type with a method, ServeHTTP, so values of that type can serve HTTP requests. Look at the implementation of the method: the receiver is a function, f, and the method calls f. That may seem odd but it's not that different from, say, the receiver being a channel and the method sending on the channel.

To make ArgServer into an HTTP server, we first modify it to have the right signature.

```go
// Argument server.
func ArgServer(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(w, os.Args)
}
```

ArgServer now has same signature as HandlerFunc, so it can be converted to that type to access its methods, just as we converted Sequence to IntSlice to access IntSlice.Sort. The code to set it up is concise:

```go
http.Handle("/args", http.HandlerFunc(ArgServer))
```

When someone visits the page /args, the handler installed at that page has value ArgServer and type HandlerFunc. The HTTP server will invoke the method ServeHTTP of that type, with ArgServer as the receiver, which will in turn call ArgServer (via the invocation f(w, req) inside HandlerFunc.ServeHTTP). The arguments will then be displayed.

In this section we have made an HTTP server from a struct, an integer, a channel, and a function, all because interfaces are just sets of methods, which can be defined for (almost) any type.

## The blank identifier in multiple assignment

If an assignment requires multiple values on the left side, but one of the values will not be used by the program, a blank identifier on the left-hand-side of the assignment avoids the need to create a dummy variable and makes it clear that the value is to be discarded.

## Import for side effect

An unused import like fmt or io in the previous example should eventually be used or removed: blank assignments identify code as a work in progress. But sometimes it is useful to import a package only for its side effects, without any explicit use.For example, during its init function, the net/http/pprof package registers HTTP handlers that provide debugging information. It has an exported API, but most clients need only the handler registration and access the data through a web page. To import the package only for its side effects, rename the package to the blank identifier:

```go
import _ "net/http/pprof"
```

This form of import makes clear that the package is being imported for its side effects, because there is no other possible use of the package: in this file, it doesn't have a name. (If it did, and we didn't use that name, the compiler would reject the program.)

## Interface checks

As we saw in the discussion of interfaces above, a type need not declare explicitly that it implements an interface. Instead, a type implements the interface just by implementing the interface's methods. In practice, most interface conversions are static and therefore checked at compile time. For example, passing an *os.File to a function expecting an io.Reader will not compile unless *os.File implements the io.Reader interface.

Some interface checks do happen at run-time, though. One instance is in the encoding/json package, which defines a Marshaler interface. When the JSON encoder receives a value that implements that interface, the encoder invokes the value's marshaling method to convert it to JSON instead of doing the standard conversion. The encoder checks this property at run time with a type assertion like:

```go
m, ok := val.(json.Marshaler)
```

If it's necessary only to ask whether a type implements an interface, without actually using the interface itself, perhaps as part of an error check, use the blank identifier to ignore the type-asserted value:

```go
if _, ok := val.(json.Marshaler); ok {
    fmt.Printf("value %v of type %T implements json.Marshaler\n", val, val)
}
```

One place this situation arises is when it is necessary to guarantee within the package implementing the type that it actually satisfies the interface. If a type—for example, json.RawMessage—needs a custom JSON representation, it should implement json.Marshaler, but there are no static conversions that would cause the compiler to verify this automatically. If the type inadvertently fails to satisfy the interface, the JSON encoder will still work, but will not use the custom implementation. To guarantee that the implementation is correct, a global declaration using the blank identifier can be used in the package:

```go
var _ json.Marshaler = (*RawMessage)(nil)
```

In this declaration, the assignment involving a conversion of a *RawMessage to a Marshaler requires that *RawMessage implements Marshaler, and that property will be checked at compile time. Should the json.Marshaler interface change, this package will no longer compile and we will be on notice that it needs to be updated.

The appearance of the blank identifier in this construct indicates that the declaration exists only for the type checking, not to create a variable. Don't do this for every type that satisfies an interface, though. By convention, such declarations are only used when there are no static conversions already present in the code, which is a rare event.
