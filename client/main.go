package main

import (
	"fmt"
	"mrpc"
)

func main() {
	rpc := mrpc.DialServer("x", "8080")

	res1 := rpc.Sum(6, 3)
	res2 := rpc.Sum(10, 20)

	res := rpc.Sum(res1, res2)
	// Access and print the fields of the struct
	fmt.Printf("Sum: %d\n", res)
	err := rpc.EndConnection()
	if err != nil {
		return
	}
}
