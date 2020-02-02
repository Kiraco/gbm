package data

import (
	"../operations"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Error struct {
	ErrorMessage string
}

// LoadData - loads all the orders and balance from the json file to a struct
func LoadData(batchFilePath string) ([]operations.Operation, []Error) {
	var errors []Error
	file, err := os.Open(batchFilePath)
	if err != nil {
		fileError := Error{}
		fileError.ErrorMessage = err.Error()
		errors = append(errors, fileError)
		return nil, errors
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var filePaths []string
	for scanner.Scan() {
		filePaths = append(filePaths, scanner.Text())
	}
	var ops []operations.Operation
	for _, path := range filePaths {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			fileError := Error{}
			fileError.ErrorMessage = err.Error()
			errors = append(errors, fileError)
			continue
		}
		operation := operations.Operation{}
		err = json.Unmarshal(file, &operation)
		if err != nil {
			fileError := Error{}
			fileError.ErrorMessage = err.Error()
			errors = append(errors, fileError)
			continue
		}
		ops = append(ops, operation)
	}
	return ops, errors
}

// PrettifyOutput - add tabs to the output for better readability
func PrettifyOutput(outputs []operations.Output) ([]string, []Error) {
	var prettyOutputs []string
	var errors []Error
	for _, output := range outputs {
		data, err := json.MarshalIndent(output, "", "\t")
		if err != nil {
			marshalError := Error{}
			marshalError.ErrorMessage = err.Error()
			errors = append(errors, marshalError)
		}
		prettyOutputs = append(prettyOutputs, string(data))
	}
	return prettyOutputs, errors
}