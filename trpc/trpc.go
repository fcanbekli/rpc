package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Find RPC functions from target directory
func findRPCFunctions(path string) ([]string, error) {
	var rpcFunctions []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(filePath) == ".go" {
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			fileContent := string(content)

			// Use a regular expression to find functions with "//RPC" on top of the definition
			rpcPattern := regexp.MustCompile(`(?s)//RPC.*?\bfunc\b[^{]+{([^}]+)}`)
			matches := rpcPattern.FindAllStringSubmatch(fileContent, -1)

			for _, match := range matches {
				if len(match) >= 2 {
					// Remove "//RPC" comment from the found function
					functionWithoutComment := strings.ReplaceAll(match[0], "//RPC", "")
					rpcFunctions = append(rpcFunctions, functionWithoutComment)
				}
			}
		}
		return nil
	})

	return rpcFunctions, err
}

func generateTrpcFile(directoryPath string) (*os.File, error) {
	// Check if the directory exists
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Directory does not exist: %s", directoryPath)
	}

	// Create a new file in the specified directory
	file, err := os.Create(filepath.Join(directoryPath, "newfile.go"))
	if err != nil {
		return nil, err
	}
	// You can write the initial content to the file if needed
	_, err = file.WriteString(InitialContent)
	_, err = file.WriteString(SerializationFunctions)
	_, err = file.WriteString(connectionContent)

	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

func addLinesAfterFunc(input string, linesToAdd string, lineNumber int) (string, error) {
	// Split the input string into lines
	lines := strings.Split(input, "\n")

	// Find the line containing "func"
	foundFunc := false
	for i, line := range lines {
		if strings.Contains(line, "func") {
			foundFunc = true

			// Insert the linesToAdd at the specified line number after "func"
			insertIndex := i + lineNumber + 1
			if insertIndex > len(lines) {
				return "", fmt.Errorf("line number %d is out of range", lineNumber)
			}

			lines = append(lines[:insertIndex], append([]string{linesToAdd}, lines[insertIndex:]...)...)
			break
		}
	}

	if !foundFunc {
		return "", fmt.Errorf("no line containing 'func' found")
	}

	// Join the modified lines back into a string
	result := strings.Join(lines, "\n")
	return result, nil
}

func generateHandleConnection(rpcFunction string) string {
	funcBody := `
func handleConnection(conn net.Conn) {
}`
	funcBody, _ = addLinesAfterFunc(funcBody, `		
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
`, 0)

	// Input serialization part
	var byteCounter = 0
	des1 := fmt.Sprintf("\ta := DeserializeInt(buffer, %d, %d)", byteCounter, byteCounter+8 /*size of int*/)
	byteCounter += 8
	funcBody, _ = addLinesAfterFunc(funcBody, des1, 10)

	des2 := fmt.Sprintf("\tb := DeserializeInt(buffer, %d, %d)", byteCounter, byteCounter+8 /*size of int*/)
	byteCounter += 8
	funcBody, _ = addLinesAfterFunc(funcBody, des2, 11)

	funcBody, _ = addLinesAfterFunc(funcBody, "result := Sum(a, b)", 12)

	str := fmt.Sprintf("\tdata := make([]byte, %d)", 8)
	funcBody, _ = addLinesAfterFunc(funcBody, str, 13)

	var byteCounter1 = 0
	ser1 := fmt.Sprintf("\tSerializeInt(data, %s, %d, %d)", "a", byteCounter1, byteCounter1+8 /*size of int*/)
	byteCounter += 8
	funcBody, _ = addLinesAfterFunc(funcBody, ser1, 14)

	funcBody, _ = addLinesAfterFunc(funcBody, `		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error sending response to client:", err)
			return
		}

		fmt.Printf("Received %d bytes from %s: %d, Sent response: %d\n", n, conn.RemoteAddr(), a, result)

		if n == 0 {
			fmt.Println("Connection closed by client")
			return
		}
	}`, 15)
	return funcBody
}

// Function to generate RPC functions with a prefix
func generateRPCFunctions(rpcFunctions []string) []string {
	var rpcWithPrefix []string

	for _, rpcFunc := range rpcFunctions {
		// Add RPC prefix to function names
		//		var finalFunction string

		finalFunction := strings.ReplaceAll(rpcFunc, "func ", "func (rpc *RpcService) RPC_")

		// Argument and return type size calculation
		argumentSize, _ := calculateArgumentSize(rpcFunc)
		//	returnSize, _ := calculateReturnSize(rpcFunc)

		// Buffer allocation
		str := fmt.Sprintf("\tdata := make([]byte, %d)", argumentSize)
		finalFunction, _ = addLinesAfterFunc(finalFunction, str, 0)

		// Input serialization part
		var byteCounter = 0
		ser1 := fmt.Sprintf("\tSerializeInt(data, %s, %d, %d)", "a", byteCounter, byteCounter+8 /*size of int*/)
		byteCounter += 8
		finalFunction, _ = addLinesAfterFunc(finalFunction, ser1, 1)

		ser2 := fmt.Sprintf("\tSerializeInt(data, %s, %d, %d)", "b", byteCounter, byteCounter+8 /*size of int*/)
		byteCounter += 8
		finalFunction, _ = addLinesAfterFunc(finalFunction, ser2, 2)

		finalFunction, _ = addLinesAfterFunc(finalFunction, "_, _ = rpc.conn.Write(data)", 3)

		// Connection read part
		finalFunction, _ = addLinesAfterFunc(finalFunction, `	// Wait for a response from the server
		responseBuffer := make([]byte, 1024)
		n, err := rpc.conn.Read(responseBuffer)
		if err != nil {
		fmt.Println("Error reading response from server:", err)
	}`, 4)

		// Deserialize response part
		var byteCounter2 = 0
		des1 := fmt.Sprintf("\tresponseValue := DeserializeInt(responseBuffer[:n], %d, %d)", byteCounter2, byteCounter2+8 /*size of int*/)
		byteCounter += 8

		finalFunction, _ = addLinesAfterFunc(finalFunction, des1, 9)

		// Return
		finalFunction, _ = addLinesAfterFunc(finalFunction, "return responseValue", 10)

		rpcWithPrefix = append(rpcWithPrefix, finalFunction)
		// Find function argument input size and return size

	}

	return rpcWithPrefix
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go /path/to/search")
		os.Exit(1)
	}

	searchPath := os.Args[1]

	rpcFunctions, err := findRPCFunctions(searchPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("RPC Functions found:")
	for _, rpcFunc := range rpcFunctions {
		fmt.Println(rpcFunc)
	}
	file, _ := generateTrpcFile(searchPath)

	x := generateRPCFunctions(rpcFunctions)

	y := generateHandleConnection(rpcFunctions[0])
	file.WriteString(y)

	for _, y := range x {
		file.WriteString(y)
	}

}
