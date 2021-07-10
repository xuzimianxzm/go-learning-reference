## Formatting

The go provides the commands program of gofmt to formatting the go code style.
which operates at the package level rather than source file level.

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

If every doc comment begins with the name of the item it describes, you can use the doc subcommand of the go tool and run the output through grep.

```shell
go doc -all regexp | grep -i item-name
```

## Package names

Notes: Don't use the import . notation, which can simplify tests that must run outside the package they are testing, but should otherwise be avoided.

## Interface names

- one-method interfaces are named by the method name plus an -er suffix or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.

## Named result parameters

The return or result "parameters" of a Go function can be given names and used as regular variables, just like the incoming parameters. When named, they are initialized to the zero values for their types when the function begins; if the function executes a return statement with no arguments, the current values of the result parameters are used as the returned values.

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

Go's defer statement schedules a function call (the deferred function) to be run immediately before the function executing the defer returns. It's an unusual but effective way to deal with situations such as resources that must be released regardless of which path a function takes to return. The canonical examples are unlocking a mutex or closing a file.

```go
// Contents returns the file's contents as a string.
func Contents(filename string) (string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()  // f.Close will run when we're finished.

    var result []byte
    buf := make([]byte, 100)
    for {
        n, err := f.Read(buf[0:])
        result = append(result, buf[0:n]...) // append is discussed later.
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err  // f will be closed if we return here.
        }
    }
    return string(result), nil // f will be closed if we return here.
}
```

Deferring a call to a function such as Close has two advantages. First, it guarantees that you will never forget to close the file, a mistake that's easy to make if you later edit the function to add a new return path. Second, it means that the close sits near the open, which is much clearer than placing it at the end of the function.

The arguments to the deferred function (which include the receiver if the function is a method) are evaluated when the defer executes, not when the call executes. Besides avoiding worries about variables changing values as the function executes, this means that a single deferred call site can defer multiple function executions. Here's a silly example.

```go
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i)
}
```

Deferred functions are executed in LIFO order, so this code will cause 4 3 2 1 0 to be printed when the function returns.

## Data

### Allocation with new

Go has two allocation primitives, the built-in functions new and make. They do different things and apply to different types, which can be confusing, but the rules are simple. Let's talk about new first. It's a built-in function that allocates memory, but unlike its namesakes in some other languages it does not initialize the memory, it only zeros it. That is, new(T) allocates zeroed storage for a new item of type T and returns its address, a value of type \*T.

- Since the memory returned by new is zeroed, it's helpful to arrange when designing your data structures that the zero value of each type can be used without further initialization.

The zero-value-is-useful property works transitively. Consider this type declaration.

```go

type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}

/*
Values of type SyncedBuffer are also ready to use immediately upon allocation or just declaration. In the next snippet, both p and v will work correctly without further arrangement.
*/
p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer

```

### Constructors and composite literals

The fields of a composite literal are laid out in order and must all be present. However, by labeling the elements explicitly as field:value pairs, the initializers can appear in any order, with the missing ones left as their respective zero values. Thus we could say.

```go
File{fd: fd, name: name}
```

- As a limiting case, if a composite literal contains no fields at all, it creates a zero value for the type. The expressions new(File) and &File{} are equivalent.

### Allocation with make

The built-in function make(T, args) serves a purpose different from new(T). It creates slices, maps, and channels only, and it returns an initialized (not zeroed) value of type T (not \*T). The reason for the distinction is that these three types represent, under the covers, references to data structures that must be initialized before use.

A slice, for example, is a three-item descriptor containing a pointer to the data (inside an array), the length, and the capacity, and until those items are initialized, the slice is nil.

These examples illustrate the difference between new and make.

```go
var p *[]int = new([]int)       // allocates slice structure; *p == nil; rarely useful
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

- We must return the slice afterwards because, although Append can modify the elements of slice, the slice itself (the run-time data structure holding the pointer, length, and capacity) is passed by value.

### Two-dimensional slices

- Go's arrays and slices are one-dimensional. To create the equivalent of a 2D array or slice, it is necessary to define an array-of-arrays or slice-of-slices, like this:

```go
type Transform [3][3]float64  // A 3x3 array, really an array of arrays.
type LinesOfText [][]byte     // A slice of byte slices.
```

### Maps

An attempt to fetch a map value with a key that is not present in the map will return the zero value for the type of the entries in the map.For instance, if the map contains integers, looking up a non-existent key will return 0.

Sometimes you need to distinguish a missing entry from a zero value. Is there an entry for "UTC" or is that 0 because it's not in the map at all? You can discriminate with a form of multiple assignment.

```go
var seconds int
var ok bool
seconds, ok = timeZone[tz]
```

For obvious reasons this is called the “comma ok” idiom. In this example, if tz is present, seconds will be set appropriately and ok will be true; if not, seconds will be set to zero and ok will be false.

## Printing

Our String method is able to call Sprintf because the print routines are fully reentrant and can be wrapped this way. There is one important detail to understand about this approach, however: don't construct a String method by calling Sprintf in a way that will recur into your String method indefinitely.

````go
type MyString string

func (m MyString) String() string {
    return fmt.Sprintf("MyString=%s", m) // Error: will recur forever.
}
````

It's also easy to fix: convert the argument to the basic string type, which does not have the method.

````go
type MyString string
func (m MyString) String() string {
    return fmt.Sprintf("MyString=%s", string(m)) // OK: note conversion.
}
````

## Append

````go
func append(slice []T, elements ...T) []T
````

where T is a placeholder for any given type. You can't actually write a function in Go where the type T is determined by the caller. That's why append is built in: it needs support from the compiler.

````go
// append element
x := []int{1,2,3}
x = append(x, 4, 5, 6)
fmt.Println(x)

// or append slice
x := []int{1,2,3}
y := []int{4,5,6}
x = append(x, y...)
fmt.Println(x)
````