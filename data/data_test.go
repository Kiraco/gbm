package data

import (
	"../operations"
	"testing"
)

func TestLoadData(t *testing.T) {
	path := []string{"/Users/donovan/Documents/Personal/Projects/go/gbm/mock-data/test.json"}
	result, errors := LoadData(path)
	if len(result) == 0 {
		t.Error("There should be a set of operations from file ")
	}
	if len(errors) != 0 {
		t.Error("There should not be any error")
	}
}

func TestLoadDataIncorrectPath(t *testing.T) {
	path := []string{"/Users/donovan/Documents/Personal/Projects/go/gbm/test"}
	result, errors := LoadData(path)

	if len(result) != 0 {
		t.Errorf("There should not be any operation")
	}
	if len(errors) == 0 {
		t.Error("There should be errors recorded")
	}
}

func TestPrettifyOutput(t *testing.T) {
	output := operations.Output{}
	result, errors := PrettifyOutput([]operations.Output{output})
	if len(result) == 0 {
		t.Error("There should be a pretty output")
	}
	if len(errors) != 0 {
		t.Error("There should not be any error")
	}
}