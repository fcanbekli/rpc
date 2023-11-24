package main

import (
	"fmt"
	"regexp"
	"strings"
)

// calculateArgumentSize calculates the total size of arguments in bytes.
func calculateArgumentSize(funcString string) (int, error) {
	// Extract argument types from the function definition
	re := regexp.MustCompile(`func\s+\w+\((.*?)\)\s+\w+\s*{`)
	matches := re.FindStringSubmatch(funcString)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid function definition")
	}

	argTypes := strings.Split(matches[1], ",")

	// Calculate the size of each argument type
	totalSize := 0
	for _, argType := range argTypes {
		argType = strings.TrimSpace(argType)
		size, err := getTypeSize(argType)
		if err != nil {
			return 0, err
		}
		totalSize += size
	}

	return totalSize, nil
}

func calculateReturnSize(funcString string) (int, error) {
	// Extract return type from the function definition
	re := regexp.MustCompile(`func\s+\w+\((.*?)\)\s*(\w+)?\s*{`)
	matches := re.FindStringSubmatch(funcString)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid function definition")
	}

	// Extract the return type, if available
	returnType := ""
	if len(matches) == 3 && matches[2] != "" {
		returnType = matches[2]
	}

	// Calculate the size of the return type
	returnType = strings.TrimSpace(returnType)
	size, err := getTypeSize(returnType)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func getTypeSize(typeString string) (int, error) {
	// Assuming int size is 8 bytes for simplicity
	if strings.Contains(typeString, "int") {
		return 8, nil
	}

	// You can extend this logic for other types as needed

	return 0, fmt.Errorf("unsupported type: %s", typeString)
}
