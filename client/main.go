package main

import (
	"fmt"
	"mrpc"
)

func main() {
	mrpc.SayHello()
	mrpc.SayMyName()
	mrpc.SayTest()

	rpc := mrpc.DialServer("x", "8080")
	
	// Access and print the fields of the struct
	fmt.Println("First Name:", rpc.Sum(2, 4))
}
