package main

import "fmt"

func main() {
	// arrays
	var a [3]int
	var b = [...]int{1, 2, 3}
	a = b
	var c = [...]int{2: 3, 1: 2}
	fmt.Println(b, a, a == b, c)

	// slices
	var d = []int{1, 2, 3}
	var e = d[1:2]
	f := make([]int, 2, 3)
	fmt.Println(d, e, f, len(f), cap(f))
}
