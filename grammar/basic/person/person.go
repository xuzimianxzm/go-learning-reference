package person

import "strconv"

type Person struct {
	ID        int //public field
	FirstName string
	LastName  string
	Address   string
	age       int // private field
}

type Employee struct {
	Person
	ManagerID int
}

type Contractor struct {
	Person
	CompanyID int
}

func (e Employee) ToString() string {
	return strconv.Itoa(e.ManagerID) + " " + strconv.Itoa(e.ID) + " " + e.FirstName + " " + e.LastName + " " + e.Address + " " + strconv.Itoa(e.age)
}

// Notice: any time we pass a value to a function, either as a receiver or as an argument, that data is copied in memory.
func (e Employee) UpdateFirstName(firstName string) {
	e.FirstName = firstName
}

/*
There are two reasons to use a pointer receiver.
1. The first is so that the method can modify the value that its receiver points to.
2. The second is to avoid copying the value on each method call. This can be more efficient if the receiver is a large struct, for example.
*/
func (e *Employee) UpdateFirstNameWithPoint(firstName string) {
	(*e).FirstName = firstName
}

func (e *Employee) UpdateLastNameWithPoint(LastName string) {
	e.LastName = LastName
}
