package utils

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestReMarshal(t *testing.T) {
	original := &TestStruct{Name: "Alice", Age: 25}
	var result TestStruct

	err := ReMarshal(original, &result)
	if err != nil {
		t.Errorf("ReMarshal() error = %v", err)
		return
	}

	if !reflect.DeepEqual(original, &result) {
		t.Errorf("ReMarshal() = %v, want %v", result, original)
	}
}
