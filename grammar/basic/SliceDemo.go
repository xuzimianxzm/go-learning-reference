package main

import "fmt"

func main() {
	mySlice := []string{"Hi", "how", "are", "you"}

	updateSliceFirst(mySlice)

	fmt.Println(mySlice)
}

/**
Go is pass value language,and when we pass a slice into a function, the function will copies the slice value.
it is still pointing at the original array in memory
*/
func updateSliceFirst(s []string) {
	s[0] = "hey"
}

func updateSliceSecond(s []string) {
	s[1] = "what"
}
