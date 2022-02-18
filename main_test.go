package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectUrl(t *testing.T) {

	_, err := getDataFromUrl(mmaUrl)

	assert.Equal(t, err, nil)
}
