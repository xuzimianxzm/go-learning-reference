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

func (e Employee) UpdateFirstName(firstName string) {
	e.FirstName = firstName
}

// UpdateFirstNameWithPoint /**
// any time we pass a value to a function, either as a receiver or as an argument, that data is copied in memory.
func (e *Employee) UpdateFirstNameWithPoint(firstName string) {
	(*e).FirstName = firstName
}

func (e *Employee) UpdateLastNameWithPoint(LastName string) {
	e.LastName = LastName
}