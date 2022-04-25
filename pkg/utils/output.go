package utils

import (
	"encoding/json"
	"fmt"
)

func OutputResult(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	fmt.Println(string(output))
	return nil
}
