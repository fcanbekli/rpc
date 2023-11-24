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

// Function to generate RPC functions with a prefix
func generateRPCFunctions(rpcFunctions []string) []string {
	var rpcWithPrefix []string

	for _, rpcFunc := range rpcFunctions {
		// Add RPC prefix to function names
		rpcWithPrefix = append(rpcWithPrefix, strings.ReplaceAll(rpcFunc, "func ", "func RPC_"))
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

	for _, y := range x {
		file.WriteString(y)
	}

}
