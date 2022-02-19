package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectUrl(t *testing.T) {

	scraper := Scraper{"https://www.sherdog.com/organizations/Ultimate-Fighting-Championship-UFC-2"}

	data, err := scraper.getDataFromUrl()

	assert.Equal(t, nil, err, "Correct URL should not return an error")
	assert.NotEqual(t, len(data), 0, "Correct URL should return some results")
}

func TestInvalidUrl(t *testing.T) {
	scraper := Scraper{"fsfjfi.f32fji"}

	_, err := scraper.getDataFromUrl()

	assert.NotEqual(t, nil, err)
}

func TestIrrelevantUrl(t *testing.T) {
	scraper := Scraper{"https://espn.co.uk"}

	_, err := scraper.getDataFromUrl()
	assert.NotEqual(t, nil, err)
}
