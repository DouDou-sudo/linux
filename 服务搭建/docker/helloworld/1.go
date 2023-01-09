package main

import "fmt"

var x int = 1

// var a string

func main1() {
	x += 1
	fmt.Println(x)
	a := "da"
	fmt.Printf("%v %T\n", a, a)
}
func main() {
	main1()
	fmt.Printf("%p %v\n", x, x)
}
