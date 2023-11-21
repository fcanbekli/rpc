package mrpc

import (
	"fmt"
	"net"
)

// Client -------------------- Client
type RpcService struct {
	conn net.Conn
}

func (rpc *RpcService) Sum(a int, b int) int {

	data := make([]byte, 8+8)

	// SERIALIZATION
	SerializeInt(data, a, 0, 8)
	SerializeInt(data, b, 8, 16)
	// SERIALIZATION

	_, _ = rpc.conn.Write(data)

	fmt.Printf("Sent message to server: %s\n", data)

	// Wait for a response from the server
	responseBuffer := make([]byte, 1024)
	n, err := rpc.conn.Read(responseBuffer)
	if err != nil {
		fmt.Println("Error reading response from server:", err)
	}

	// SERIALIZATION
	responseValue := DeserializeInt(responseBuffer[:n], 0, 8)
	// SERIALIZATION

	return responseValue
}

func (rpc *RpcService) EndConnection() error {
	return rpc.conn.Close()
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

// Client -------------------- Client

// Server -------------------- Server
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

		// SERIALIZATION
		a := DeserializeInt(buffer, 0, 8)
		b := DeserializeInt(buffer, 8, 16)
		// SERIALIZATION

		result := sum_rpc(a, b)

		data := make([]byte, 8)
		SerializeInt(data, result, 0, 8) //TODO: Make a version of this which is only serialize single integer
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

func sum_rpc(a int, b int) int {
	return a + b
	// Real calculation here
}
