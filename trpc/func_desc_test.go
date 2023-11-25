package main

import (
	"testing"
)

func TestFcGetArgumentSize(t *testing.T) {
	desc := &FunctionDescription{}

	desc.ArgumentTypes = []string{"int"}
	size := desc.GetArgumentSize()
	if size != 8 {
		t.Errorf("Expected size %d, got %d", 8, size)
	}

	desc.ArgumentTypes = []string{"int", "int"}
	size = desc.GetArgumentSize()
	if size != 16 {
		t.Errorf("Expected size %d, got %d", 16, size)
	}
}

func TestFcGetReturnTypeSize(t *testing.T) {
	desc := &FunctionDescription{}

	desc.ReturnTypes = []string{"int"}
	size := desc.GetReturnTypeSize()
	if size != 8 {
		t.Errorf("Expected size %d, got %d", 8, size)
	}

	desc.ReturnTypes = []string{"float32", "float32"}
	size = desc.GetReturnTypeSize()
	if size != 8 {
		t.Errorf("Expected size %d, got %d", 8, size)
	}
}

func TestFcNew(t *testing.T) {
	tests := []struct {
		rawFunction      string
		expectedName     string
		expectedArgTypes []string
		expectedRetTypes []string
	}{ //TODO: Add more test cases to here!
		{
			rawFunction: `
				// RPC
				func Sum(a int, b int) int {
					return a + b
				}
			`,
			expectedName:     "Sum",
			expectedArgTypes: []string{"int", "int"},
			expectedRetTypes: []string{"int"},
		},
		// Add more test cases here
	}

	for _, test := range tests {
		desc := &FunctionDescription{}
		desc.New(test.rawFunction)

		// Check function name
		if desc.FunctionName != test.expectedName {
			t.Errorf("Expected function name %s, got %s", test.expectedName, desc.FunctionName)
		}

		// Check argument types
		if !equalStringSlices(desc.ArgumentTypes, test.expectedArgTypes) {
			t.Errorf("Expected argument types %v, got %v", test.expectedArgTypes, desc.ArgumentTypes)
		}

		// Check return types
		if !equalStringSlices(desc.ReturnTypes, test.expectedRetTypes) {
			t.Errorf("Expected return types %v, got %v", test.expectedRetTypes, desc.ReturnTypes)
		}
	}
}

// Helper function to check equality of string slices
func equalStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}
	return true
}
