package main

import (
	"log"
	"xuzimian.com/grammar/basic/person"
)

func main() {
	var employ person.Employee
	var person person.Person
	// 下面这一行会报错,因为变量名和包名冲突了
	// var employ person.Employee

	log.Print("hello", person, employ)
}
