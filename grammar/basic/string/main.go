package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	f := "Golang编程"

	// 统计的是byte长度，它的类型为 'unit8',代表了一个ASCII字符。由于中文字符在UTF-8中占用了3个字节，所以使用len方法时获得的中文长度是6个字节。
	fmt.Printf("byte len of f is %v\n", len(f))
	// 统计的是rune类型，它的类型为 'int32',代表了一个UTF-8字符，它可以类比为Java中的char类型。
	// utf8.RuneCountInString()方法统计的是字符串的 Unicode字符数量。
	fmt.Printf("rune len of f is %v\n", utf8.RuneCountInString(f))

	// 在进行字节遍历时，因为中文字符的Unicode字符会被阶段，导致中文输出乱码
	for _, g := range []byte(f) {
		fmt.Printf("%c", g)
	}

	fmt.Println()
	for _, h := range f {
		fmt.Printf("%c", h)
	}
}
