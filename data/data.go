package data

import (
	"../operations"
	"encoding/json"
	"io/ioutil"
)

func LoadData(filePath string) operations.Operation {
	file, _ := ioutil.ReadFile(filePath)
	operation := operations.Operation{}
	err := json.Unmarshal(file, &operation)
	if err != nil {
		panic(err)
	}
	return operation
}

func PrettifyOutput(output operations.Output) ([]byte, error){
	return json.MarshalIndent(output, "", "\t")
}