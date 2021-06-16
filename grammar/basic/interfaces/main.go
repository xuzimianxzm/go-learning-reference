package main

import (
	"fmt"
	"math"
)

type I interface {
	M()
}

type T struct {
	S string
}

// M : this method is T struct implemented the interface of I, because interfaces are implemented implicitly,and
// there is no explicit declaration of intent, no "implements" keyword.
func (t *T) M() {
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

type F float64

func (f F) M() {
	fmt.Println(f)
}

//One of the most ubiquitous interfaces is Stringer defined by the fmt package.
//A Stringer is a type that can describe itself as a string. The fmt package (and many others) look for this interface to print values.
func (f F) String() string {
	return fmt.Sprintf("number is %d", f)
}

func main() {
	var i I

	var t *T
	i = t
	describe(i)
	i.M()

	i = F(math.Pi)
	describe(i)
	i.M()

	i = &T{"hello"}
	describe(i)
	i.M()

}

// Interface values
// Under the hood, interface values can be thought of as a tuple of a value and a concrete type:
// (value, type)
// Calling a method on an interface value executes the method of the same name on its underlying type.
func describe(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}
