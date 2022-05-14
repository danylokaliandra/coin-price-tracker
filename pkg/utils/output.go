package utils

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func OutputResult(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal output: %s", err.Error())
	}

	log.Info(string(output))

	return nil
}
