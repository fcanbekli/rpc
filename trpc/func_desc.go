package main

import (
	"regexp"
	"strings"
	"unsafe"
)

type FunctionDescription struct {
	FunctionName  string
	FunctionRaw   string
	ArgumentTypes []string // int, float32, etc...
	ReturnTypes   []string // int, float32, etc...
}

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
		}
	}

	return fc
}

func (fc *FunctionDescription) GetArgumentSize() int {
	var totalSize int
	for _, s := range fc.ArgumentTypes {
		totalSize += getSizeMap()[s]
	}
	return totalSize
}

func (fc *FunctionDescription) GetReturnTypeSize() int {
	var totalSize int
	for _, s := range fc.ReturnTypes {
		totalSize += getSizeMap()[s]
	}
	return totalSize
}

func getSizeMap() map[string]int {
	sizeMap := make(map[string]int)

	sizeMap["int"] = int(unsafe.Sizeof(int(0)))
	sizeMap["int8"] = int(unsafe.Sizeof(int8(0)))
	sizeMap["int16"] = int(unsafe.Sizeof(int16(0)))
	sizeMap["int32"] = int(unsafe.Sizeof(int32(0)))
	sizeMap["int64"] = int(unsafe.Sizeof(int64(0)))

	sizeMap["uint"] = int(unsafe.Sizeof(uint(0)))
	sizeMap["uint8"] = int(unsafe.Sizeof(uint8(0)))
	sizeMap["uint16"] = int(unsafe.Sizeof(uint16(0)))
	sizeMap["uint32"] = int(unsafe.Sizeof(uint32(0)))
	sizeMap["uint64"] = int(unsafe.Sizeof(uint64(0)))

	sizeMap["float32"] = int(unsafe.Sizeof(float32(0)))
	sizeMap["float64"] = int(unsafe.Sizeof(float64(0)))

	sizeMap["bool"] = int(unsafe.Sizeof(false))

	sizeMap["string"] = int(unsafe.Sizeof(""))
	return sizeMap
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
