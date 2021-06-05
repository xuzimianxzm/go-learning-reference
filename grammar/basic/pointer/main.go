package main

import (
	"fmt"
	. "xuzimian.com/grammar/basic/person"
)

/*  '*variable' 语法是将内存地址转换成该内存地址中存放的值
 *  '&variable' 语法是将变量的值转换成改变量的内存地址
 */
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
