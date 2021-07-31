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
- 切片的结构和字符串类似，但是解除了制度限制。切片可以和nil比较，只有当切片底层数组为空时切片本身才为nil,这时候切片的len和cap的信息是无效的，如
  果有切片底层数据指针为空但len和cap都不为0的情况，那么说明切片本身已经损坏了(reflect.SliceHeader或unsafe包对切片做了不正确的修改)。
- 空数组一般很少用到，但对于切片来说，len为0但cap容量不为0的切片是非常有用的特性。例如下面的TrimSpace()函数用于删除[]byte中的空格，函数利用了 长度为0的切片特性，实现高效且简洁:
  ```go
  func TrimSpace(s []byte) []byte{
    b := s[:0]
    for _,x := range s {
      if x != ' ' {
       b = append(b, x)
      }
    }
    return b
  }
  ```
- Go语言字符串底层数据对应的也是字节数组，但是字符串的制度属性禁止了在程序中对地秤字节数组的元素的修改。字符串赋值只是复制了数据地址和对应的长度，
  而不会导致底层数据的复制。字符串虽然不是切片，但是支持切片操作。不同位置的切片底层访问的是同一块内存数据。
- Go语言的源文件都采用UTF8编码。因此，Go源文件中出现的字符串面值常量一般也是UTF8编码的(对于转译字符则没有这个限制)，一般都假设Go字符串对应的是 一个合法的UTF8编码的字符序列，可以用for
  range循环直接偏离UTF8解码后的Unicode码点值。
  
  

