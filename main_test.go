package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectUrl(t *testing.T) {

	scraper := Scraper{mmaUrl}

	_, err := scraper.getDataFromUrl()

	assert.Equal(t, nil, err)
}
