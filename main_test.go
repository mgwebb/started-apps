package main

import (
        "testing"
       )

func TestConvertSizeG(t *testing.T) {
	actualResult := ConvertSize(2048)
	var expectedResult = "2G"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestConvertSizeM(t *testing.T) {
	actualResult := ConvertSize(768)
	var expectedResult = "768M"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
