package mrpc

import (
	"encoding/binary"
	"fmt"
	"net"
)

// serializeInts serializes two integers into a byte slice
func serializeInts(a, b int) []byte {
	// Create a byte slice with enough capacity to hold two integers
	data := make([]byte, 8*2)

	// Use encoding/binary to serialize integers into the byte slice
	binary.BigEndian.PutUint64(data[0:8], uint64(a))
	binary.BigEndian.PutUint64(data[8:16], uint64(b))

	return data
}

// deserializeInts deserializes two integers from a byte slice
func deserializeInts(data []byte) (int, int) {
	// Use encoding/binary to deserialize integers from the byte slice
	a := int(binary.BigEndian.Uint64(data[0:8]))
	b := int(binary.BigEndian.Uint64(data[8:16]))

	return a, b
}

// Client -------------------- Client
type RpcService struct {
	conn net.Conn
}

func (rpc *RpcService) Sum(a int, b int) int {

	// Send a "Hello, World!" message to the server
	message := serializeInts(a, b)
	_, _ = rpc.conn.Write(message)

	fmt.Printf("Sent message to server: %s\n", message)

	// Wait for a response from the server
	responseBuffer := make([]byte, 1024)
	n, err := rpc.conn.Read(responseBuffer)
	if err != nil {
		fmt.Println("Error reading response from server:", err)
	}

	// Deserialize the response
	responseValue, _ := deserializeInts(responseBuffer[:n])

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

	// Create a buffer to read data from the connection
	buffer := make([]byte, 1024)

	for {
		// Read data from the connection
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// Print the data received from the connection
		a, b := deserializeInts(buffer)
		result := sum_rpc(a, b)

		// Send the result back to the client
		response := serializeInts(result, 0) //TODO: Make a version of this which is only serialize single integer
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error sending response to client:", err)
			return
		}

		// Print the received data and result
		fmt.Printf("Received %d bytes from %s: %d, Sent response: %d\n", n, conn.RemoteAddr(), a, result)

		// Check if the connection is closed by the client
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
