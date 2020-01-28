package main

import (
	"./data"
	"./operations"
	"fmt"
	"os"
)

func main() {
	filepath := os.Args[1]
	operation := data.LoadData(filepath)
	result := operations.PerformOperation(&operation)
	output, error :=data.PrettifyOutput(result)
	if error != nil {
		panic("error, couldn't parse output. Something bad happened.")
	}
	fmt.Printf("%+v\n", string(output))
}
