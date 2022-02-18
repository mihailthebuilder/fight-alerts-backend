package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectUrl(t *testing.T) {

	received := getDataFromUrl(mmaUrl)

	fmt.Println(received, nil)

	assert.NotEqual(t, received, nil)
}
