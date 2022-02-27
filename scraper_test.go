package main

import (
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

	today := time.Now()
	future := today.Add(time.Hour * 24 * 10)

	return future.String()
}

func (e MockCollyElem) ChildText(selector string) string {
	if !e.validHeadline {
		return ""
	}

	return "valid headline"
}

func ValidateFightRecord(record FightRecord, t *testing.T) {
	if len([]rune(record.Headline)) == 0 || time.Now().After(record.DateTime) {
		t.Errorf("record should not have nil or nil-equivalent values - %#v", record)
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
			args:    args{e: MockCollyElem{validDate: false, validHeadline: true}},
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
