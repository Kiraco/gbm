package data

import (
	"../operations"
	"encoding/json"
	"io/ioutil"
)

type Error struct {
	ErrorMessage string
}

func LoadData(filePaths []string) ([]operations.Operation, []Error) {
	var ops []operations.Operation
	var errors []Error
	for _, path := range filePaths {
		file, _ := ioutil.ReadFile(path)
		operation := operations.Operation{}
		err := json.Unmarshal(file, &operation)
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