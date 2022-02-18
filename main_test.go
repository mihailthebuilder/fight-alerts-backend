package main

import (
	"testing"
)

func TestCorrectUrl(t *testing.T) {

	if getDataFromUrl(mmaUrl) == nil {
		t.Error("Test with correct URL failed")
	}
}
