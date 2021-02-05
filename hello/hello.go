package main

import (
	"fmt"

	"xuzimian.com/greetings"
)

func main() {
	message := greetings.Hello("go")
	fmt.Println(message)
}
