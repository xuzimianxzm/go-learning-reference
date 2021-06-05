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
	// &variable 返回这个变量所指向的内存地址， *variable 返回这个内存地址所指向的值
	fmt.Println(&employeePointer, " ", *employeePointer)

	employee.UpdateFirstNameWithPoint("another")
	fmt.Println("After call UpdateFirstNameWithPoint,point:", employee.ToString())
	fmt.Println("After call UpdateFirstNameWithPoint,employ:", employeePointer.ToString())
}
