package utils

import (
	"encoding/json"
	"fmt"
)

func OutputResult(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal output: %s", err.Error())
	}

	fmt.Println(string(output))

	return nil
}
