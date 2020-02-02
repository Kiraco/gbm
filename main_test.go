package main

import (
	"./data"
	"testing"
)

func TestRunBatch(t *testing.T) {
	operations, _ := data.LoadData([]string{"/Users/donovan/Documents/Personal/Projects/go/gbm/mock-data/test.json"})
	output := RunBatch(&operations)
	if len(output) == 0 {
		t.Error("There should be an output for the operation")
	}
}

func TestMainProgram(m *testing.T) {
	// Not needed in coverage, since is the run per se
}
