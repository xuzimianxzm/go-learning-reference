package main

import "time"

func main() {
	ticker := time.NewTicker(1000)
	defer ticker.Stop()

	for i := 0; i <= 100; i++ {
		select {
		case chanTime := <-ticker.C:
			println(chanTime.String())
		default:
			println("------")
		}
	}

	n := 0
	var dur time.Duration = time.Microsecond
	chRate := time.Tick(dur)
	for chanTime := range chRate {
		println(chanTime.String())
		if n > 100 {
			return
		}
		n++
	}
}
