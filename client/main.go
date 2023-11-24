package main

import (
	"fmt"
)

func main() {
	rpc := DialServer("x", "8080")

	res := rpc.RPC_Sum(10, 3)
	// Access and print the fields of the struct
	fmt.Printf("Sum: %d\n", res)
	//err := rpc.EndConnection()
}
