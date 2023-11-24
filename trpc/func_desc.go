package main

import (
	"regexp"
	"strings"
)

type FunctionDescription struct {
	FunctionName               string
	FunctionRaw                string
	ArgumentTypes              []string // int, float32, etc...
	ArgumentsTotalSizeAsByte   int
	ReturnTypes                []string // int, float32, etc...
	ReturnTypesTotalSizeAsByte int
}

// New
func (fc *FunctionDescription) New(functionRaw string) *FunctionDescription { //TODO: This code is does not working as expected, it only works for Sum(a int, b int) int case
	fc.FunctionRaw = functionRaw

	// Extract function name
	nameRegex := regexp.MustCompile(`func\s+([a-zA-Z_]\w*)\s*\(`)
	matches := nameRegex.FindStringSubmatch(functionRaw)
	if len(matches) > 1 {
		fc.FunctionName = matches[1]
	}

	// Extract argument types and calculate total size
	argRegex := regexp.MustCompile(`\(([^)]*)\)`)
	argMatches := argRegex.FindStringSubmatch(functionRaw)
	if len(argMatches) > 1 {
		argTypes := strings.Split(argMatches[1], ",")
		fc.ArgumentTypes = make([]string, len(argTypes))
		for i, argType := range argTypes {
			// Extract only the type part (ignore variable names)
			parts := strings.Fields(argType)
			if len(parts) > 0 {
				fc.ArgumentTypes[i] = parts[len(parts)-1] // Take the last part as type
				fc.ArgumentsTotalSizeAsByte += typeSizeInBytes(parts[len(parts)-1])
			}
		}
	}

	// Extract return types and calculate total size
	retRegex := regexp.MustCompile(`\)\s*([^)]*)\s*\{`)
	retMatches := retRegex.FindStringSubmatch(functionRaw)
	if len(retMatches) > 1 {
		retTypes := strings.Split(retMatches[1], ",")
		fc.ReturnTypes = make([]string, len(retTypes))
		for i, retType := range retTypes {
			fc.ReturnTypes[i] = strings.TrimSpace(retType)
			fc.ReturnTypesTotalSizeAsByte += typeSizeInBytes(strings.TrimSpace(retType))
		}
	}

	return fc
}

// Helper function to calculate size of types in bytes
func typeSizeInBytes(typeName string) int {
	switch typeName {
	case "int", "int8", "uint8", "byte":
		return 8
	case "int16", "uint16":
		return 2
	case "int32", "uint32", "float32":
		return 4
	case "int64", "uint64", "float64":
		return 8
	default:
		// For simplicity, other types are considered to be of size 0
		return 0
	}
}
