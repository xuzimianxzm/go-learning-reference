package main

import (
	"fmt"
	"time"
)

/* [无缓存的]Channel上的发送操作,总在对应的接收操作 [完成前] 发生;对于从[无缓存]Channel进行的接收，发生在对该Channel进行的发送 [完成之前]。
即，只要是无缓存的channel,如果先发生的是对其进行接收，则会block,直到对该channel完成了发送操作。如果先发生的是对无缓存的channel进行发送，
则会block，直到对该channel的接收操作完成。
*/
func main() {
	rule1()
}

// [无缓存的]Channel上的发送操作,总在对应的接收操作 [完成前] 发生
func rule1() {
	done := make(chan int)

	go func() {
		fmt.Println("对 channel 发送前")
		time.Sleep(time.Second * 5)
		done <- 1
		fmt.Println("对 channel 发送后")
	}()

	fmt.Println("对 channel 接收前")
	<-done
	fmt.Println("对 channel 接收后")
}

// 对于从[无缓存]Channel进行的接收，发生在对该Channel进行的发送 [完成之前]。
func rule2() {
	done := make(chan int)

	go func() {
		fmt.Println("对 channel 接收前")
		<-done
		fmt.Println("对 channel 接收后")
	}()

	fmt.Println("对 channel 发送前")
	time.Sleep(time.Second * 5)
	done <- 1
	fmt.Println("对 channel 发送后")
}
