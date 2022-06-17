package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type DateRange struct {
	StartDate int64
	EndDate   int64
}

// Returns a list of Unix timestamps range to fetch a certain amount in time of historical data
func GetDateRanges(maxHistory int, timeRange int) []DateRange {
	now := time.Now()
	dateRanges := []DateRange{}

	for d := now; d.Unix() >= now.AddDate(0, 0, -maxHistory).Unix(); {
		dateRanges = append(dateRanges, DateRange{
			StartDate: d.AddDate(0, 0, -15).UnixMilli(),
			EndDate:   d.UnixMilli(),
		})
		d = d.AddDate(0, 0, -timeRange)
	}

	return dateRanges
}

// Return complete path where data is stored
func GetDataPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user homedir: %w", err)
	}

	return fmt.Sprintf("%s/.tracklet/data", homeDir), nil
}

// Load json data from file
func LoadFromFile(filename string) []byte {
	dataPath, err := GetDataPath()
	if err != nil {
		log.Errorf("could not get data path: %v", err)
		return nil
	}

	filename = fmt.Sprintf("%s/%s", dataPath, filename)

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Errorf("Could not open file: %v", err)
		return nil
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("Could not read data: %v", err)
		return nil
	}

	return data
}

// Write json data to file
func WriteToFile(filename string, content interface{}) error {
	dataPath, err := GetDataPath()
	if err != nil {
		return fmt.Errorf("could not get data path: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s.json", dataPath, filename)

	log.Infof("Saving data to file '%s'", dataPath)

	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal output: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("could not write data to file: %w", err)
	}

	return nil
}
