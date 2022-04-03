package scraper

import (
	"reflect"
	"testing"
	"time"
)

type MockCollyElem struct {
	validDate     bool
	validHeadline bool
}

func (e MockCollyElem) ChildAttr(selector string, attr string) string {
	if !e.validDate {
		return "invalid date"
	}

	return "2050-03-05T00:00:00-08:00"
}

func (e MockCollyElem) ChildText(selector string) string {
	if !e.validHeadline {
		return ""
	}

	return "valid headline"
}

func ValidateFightRecord(record FightRecord, t *testing.T) {
	if len([]rune(record.Headline)) == 0 {
		t.Errorf("record should have headline - %#v", record)
	}

	currentDayStart := time.Now().UTC().Truncate(time.Hour * 24)
	fightDayStart := record.DateTime.UTC().Truncate(time.Hour * 24)

	if fightDayStart.Before(currentDayStart) {
		t.Errorf("fight date should be in the future - %#v", record)
	}
}

func Test_parseCollyHtml(t *testing.T) {
	type args struct {
		e ICollyHtmlElem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid date",
			args:    args{e: MockCollyElem{validDate: false, validHeadline: true}},
			wantErr: true,
		},
		{
			name:    "invalid headline",
			args:    args{e: MockCollyElem{validDate: true, validHeadline: false}},
			wantErr: true,
		},
		{
			name:    "valid record",
			args:    args{e: MockCollyElem{validDate: true, validHeadline: true}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCollyHtml(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCollyHtml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				ValidateFightRecord(got, t)
			}
		})
	}
}

func Test_parseDateTime(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name     string
		args     args
		wantYear int
		wantErr  bool
	}{
		{
			name:     "valid date",
			args:     args{s: "2022-03-05T00:00:00-08:00"},
			wantYear: 2022,
			wantErr:  false,
		},
		{
			name:     "invalid date",
			args:     args{s: "invalid date"},
			wantYear: 1,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateTime(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Year(), tt.wantYear) {
				t.Errorf("parseDateTime() = %v / year %v, want %v", got, got.Year(), tt.wantYear)
			}
		})
	}
}

func Test_getDataFromUrl(t *testing.T) {
	var tests = []struct {
		input       string
		wantResults bool
		wantError   bool
	}{{MmaUrl, true, false}, {"ffefw.fdfsfs", false, true}, {"https://espn.co.uk", false, true}}

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
