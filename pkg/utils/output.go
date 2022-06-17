package utils

import (
	"encoding/json"
	"fmt"
)

// Print json structures
func OutputResult(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal output: %w", err)
	}

	fmt.Println(string(output))

	return nil
}
