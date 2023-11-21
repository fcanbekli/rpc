package main

import (
	"fmt"
	"mrpc"
)

func main() {
	rpc := mrpc.DialServer("x", "8080")

	res := rpc.Sum(1, 3)
	// Access and print the fields of the struct
	fmt.Printf("Sum: %d\n", res)
	err := rpc.EndConnection()
	if err != nil {
		return
	}
}
