package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectUrl(t *testing.T) {

	received, err := getDataFromUrl(mmaUrl)

	assert.Equal(t, err, nil)
	assert.NotEqual(t, received, nil)
}
