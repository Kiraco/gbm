package data

import (
	"../operations"
	"encoding/json"
	"io/ioutil"
)

func LoadData(filePaths []string) []operations.Operation {
	var ops []operations.Operation
	for _, path := range filePaths {
		file, _ := ioutil.ReadFile(path)
		operation := operations.Operation{}
		err := json.Unmarshal(file, &operation)
		if err != nil {
			panic(err)
		}
		ops = append(ops, operation)
	}
	return ops
}

func PrettifyOutput(outputs []operations.Output) []string {
	var prettyOutputs []string
	for _, output := range outputs {
		data, err := json.MarshalIndent(output, "", "\t")
		if err != nil {
			panic(err)
		}
		prettyOutputs = append(prettyOutputs, string(data))
	}
	return prettyOutputs
}