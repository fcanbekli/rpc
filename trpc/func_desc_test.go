package main

import (
	"testing"
)

func TestFunctionDescriptionNew(t *testing.T) {
	tests := []struct {
		rawFunction      string
		expectedName     string
		expectedArgTypes []string
		expectedArgsSize int
		expectedRetTypes []string
		expectedRetSize  int
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
			expectedArgsSize: 16, // (int: 4 bytes) + (int: 4 bytes) = 8 bytes
			expectedRetTypes: []string{"int"},
			expectedRetSize:  8, // (int: 4 bytes)
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

		// Check arguments total size
		if desc.ArgumentsTotalSizeAsByte != test.expectedArgsSize {
			t.Errorf("Expected arguments total size %d bytes, got %d bytes", test.expectedArgsSize, desc.ArgumentsTotalSizeAsByte)
		}

		// Check return types
		if !equalStringSlices(desc.ReturnTypes, test.expectedRetTypes) {
			t.Errorf("Expected return types %v, got %v", test.expectedRetTypes, desc.ReturnTypes)
		}

		// Check return types total size
		if desc.ReturnTypesTotalSizeAsByte != test.expectedRetSize {
			t.Errorf("Expected return types total size %d bytes, got %d bytes", test.expectedRetSize, desc.ReturnTypesTotalSizeAsByte)
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
