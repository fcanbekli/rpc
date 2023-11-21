package main

import (
	"fmt"
	"mrpc"
)

func main() {
	fmt.Println("Tiny RPC Server Starting")

	mrpc.StartServer("8080")
}
