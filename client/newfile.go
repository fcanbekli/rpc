// This file is auto generated, do not change content
package main

import (
	"fmt"
	"net"
	"encoding/binary"
)

func SerializeInt(data []byte, value int, left int, right int) {
	// Use encoding/binary to serialize the integer into the byte slice
	binary.BigEndian.PutUint64(data[left:right], uint64(value))
}

func DeserializeInt(data []byte, start int, end int) int {
	// Use encoding/binary to deserialize the integer from the byte slice
	return int(binary.BigEndian.Uint64(data[start:end]))
}

type RpcService struct {
	conn net.Conn
}

func DialServer(ip string, port string) *RpcService {

	conn, err := net.Dial("tcp", "localhost:8080")
	rpc_service := RpcService{}
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	rpc_service.conn = conn

	return &rpc_service
}


func StartServer(port string) {
	//Start TCP Socket
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
		
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

	a := DeserializeInt(buffer, 0, 8)
	b := DeserializeInt(buffer, 8, 16)
result := Sum(a, b)
	data := make([]byte, 8)
	SerializeInt(data, a, 0, 8)
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending response to client:", err)
			return
		}

		fmt.Printf("Received %d bytes from %s: %d, Sent response: %d\n", n, conn.RemoteAddr(), a, result)

		if n == 0 {
			fmt.Println("Connection closed by client")
			return
		}
	}
}
func (rpc *RpcService) RPC_Sum(a int, b int) int {
	data := make([]byte, 16)
	SerializeInt(data, a, 0, 8)
	SerializeInt(data, b, 8, 16)
_, _ = rpc.conn.Write(data)
	// Wait for a response from the server
		responseBuffer := make([]byte, 1024)
		n, err := rpc.conn.Read(responseBuffer)
		if err != nil {
		fmt.Println("Error reading response from server:", err)
	responseValue := DeserializeInt(responseBuffer[:n], 0, 8)
return responseValue
	}
	return a + b
}
func (rpc *RpcService) RPC_Multiply(a int, b int) int {
	data := make([]byte, 16)
	SerializeInt(data, a, 0, 8)
	SerializeInt(data, b, 8, 16)
_, _ = rpc.conn.Write(data)
	// Wait for a response from the server
		responseBuffer := make([]byte, 1024)
		n, err := rpc.conn.Read(responseBuffer)
		if err != nil {
		fmt.Println("Error reading response from server:", err)
	responseValue := DeserializeInt(responseBuffer[:n], 0, 8)
return responseValue
	}
	return a + b
}