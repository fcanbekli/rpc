package mrpc

import "fmt"

// SayHello prints a greeting message.
func SayHello() {
	fmt.Println("DATA")
}

func SayMyName() {
	fmt.Println("MyName")
}

func SayTest() {
	fmt.Println("Test")
}

// Client -------------------- Client
type RpcService struct{}

func (rpc *RpcService) Sum(a int, b int) int {
	return sum_rpc(a, b)
	// Send method name and data via tcp connection
	// Wait until have response
}

func DialServer(ip string, port string) RpcService {
	return RpcService{}
}

// Client -------------------- Client

// Server -------------------- Server
func StartServer(port string) {
	//Start TCP Socket
}

func sum_rpc(a int, b int) int {
	return a + b
	// Real calculation here
}
