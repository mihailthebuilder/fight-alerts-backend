package scraper

import (
	"fight-alerts-backend/resources"
	"testing"
)

func Test_getDataFromUrl(t *testing.T) {
	var tests = []struct {
		input       string
		wantResults bool
		wantError   bool
	}{{resources.MmaUrl, true, false}, {"ffefw.fdfsfs", false, true}, {"https://espn.co.uk", false, true}}

	for _, test := range tests {
		var scraper IScraper = Scraper{test.input}

		results, err := scraper.GetResultsFromUrl()
		gotResults, gotError := len(results) > 0, err != nil

		if gotError != test.wantError {
			t.Errorf("getDataFromUrl(%v) error = %#v | want error = %v", test.input, err.Error(), test.wantError)
		}

		if gotResults != test.wantResults {
			t.Errorf("getDataFromUrl(%v) results = %#v | want results = %v", test.input, results, test.wantResults)
		}

		if gotResults {
			for _, record := range results {
				ValidateFightRecord(record, t)
			}
		}
	}
}
