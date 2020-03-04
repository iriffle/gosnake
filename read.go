package main

import (
	"github.com/google/logger"
	"io/ioutil"
	"os"
)

// readData reads game data from a JSON file and returns the byteData
func ReadFile(file string) []byte {
	// Check if high score file exists. If not then create it
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		_, err := os.Create(file)
		if err != nil {
			logger.Errorf("Error creating file: %v", err)
		}
	}
	logger.Info("File exists.")

	// Open the data file
	jsonFile, err := os.Open(file)
	if err != nil {
		logger.Errorf("Error opening file: %v", err)
	}
	logger.Infof("Opened file for reading: %v", file)

	// Close the data file on exit
	defer func() {
		if err = jsonFile.Close(); err != nil {
			logger.Errorf("Error closing file: %v", err)
		}
		logger.Infof("Closed file after reading: %v", file)
	}()

	// Read bytes of data from JSON file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}
