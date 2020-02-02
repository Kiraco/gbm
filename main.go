package main

import (
	"./data"
	"./operations"
	"fmt"
	"os"
	"sync"
)

func main() {
	filepath := os.Args[1:]
	allFiles, _ := data.LoadData(filepath)
	result := RunBatch(&allFiles)
	outputs, _ := data.PrettifyOutput(result)
	fmt.Printf("%+v\n", outputs)
}

func RunBatch(allFiles *[]operations.Operation) []operations.Output {
	var outputs []operations.Output
	var wg sync.WaitGroup
	for _, file := range *allFiles {
		wg.Add(1)
		outputs = append(outputs, operations.PerformOperation(&file, &wg))
	}
	wg.Wait()
	return outputs
}
