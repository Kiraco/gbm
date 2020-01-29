package main

import (
	"./data"
	"./operations"
	"fmt"
	"os"
)

func main() {
	filepath := os.Args[1:]
	operation := data.LoadData(filepath)
	result := operations.RunBatch(&operation)
	outputs := data.PrettifyOutput(result)
	fmt.Printf("%+v\n", outputs)
}
