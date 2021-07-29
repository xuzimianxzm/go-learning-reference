package main

import "fmt"

type IReader interface {
	Read(p string) int
}

type IWriter interface {
	Write(p string) int
}

type IReadWriter interface {
	IReader
	IWriter
}

type Reader struct {
	buf []string
}

type Writer struct {
}

func (r Reader) Read(p string) int {
	r.buf = append(r.buf, p)
	return len(r.buf)
}

func (w Writer) Write(p string) int {
	return len(p)
}

type ReadWriter struct {
	*Reader // *bufio.Reader
	*Writer // *bufio.Writer
}

// 手动实现ReadWriter的IReader接口,可以通过ReadWriter调用Read方法，
func (rw *ReadWriter) Read(p string) int {
	return rw.Reader.Read(p)
}

func main() {
	rw :=ReadWriter{ &Reader{[]string{} },&Writer{}}
	fmt.Println(rw.Read("hello world"))
	// 未通过ReadWriter实现IWriter的接口，但可以直接调用。
	fmt.Println(rw.Write("hello world"))
}
