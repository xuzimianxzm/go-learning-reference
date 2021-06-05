package main

import (
	"fmt"
	. "xuzimian.com/grammar/basic/person"
)

func main() {
	employee := Employee{ManagerID: 1,
		Person: Person{ID: 1, FirstName: "xu", LastName: "zimian", Address: "wuhan"},
	}

	fmt.Println("Before call UpdateFirstName:", employee.ToString())
	employee.UpdateFirstName("any")
	fmt.Println("After call UpdateFirstName:", employee.ToString())

	employeePointer := &employee
	employeePointer.UpdateFirstNameWithPoint("any")
	fmt.Println("After call UpdateFirstNameWithPoint,point:", employee.ToString())
	fmt.Println("After call UpdateFirstNameWithPoint,employ:", employeePointer.ToString())
}
