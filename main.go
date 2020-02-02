package main

import (
	"./data"
	"./operations"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	filepath := os.Args[1]
	allFiles, errors := data.LoadData(filepath)
	result := RunBatch(&allFiles)
	outputs, _ := data.PrettifyOutput(result)
	for i, out := range outputs {
		output := fmt.Sprintf("%+v", out)
		fileName := fmt.Sprintf("gbm_file_%d_%d.json", i, time.Now().Unix())
		file, _ := os.Create(fileName)
		defer file.Close()
		file.WriteString(output)
	}
	if len(errors) > 0 {
		fileNameErrors := fmt.Sprintf("gbm_output_errors_%d.json", time.Now().Unix())
		fileErrors, _ := os.Create(fileNameErrors)
		defer fileErrors.Close()
		errorsOutput := fmt.Sprintf("%+v\n", errors)
		fileErrors.WriteString(errorsOutput)
	}
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
