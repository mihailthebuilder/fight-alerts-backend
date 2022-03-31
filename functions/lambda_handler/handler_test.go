package handler

import (
	"fight-alerts-backend/scraper"
	"fmt"
	"testing"
)

type MockScraper struct {
	returnError bool
}

type MockDatastore struct {
	returnError bool
}

func (s MockScraper) GetResultsFromUrl() ([]scraper.FightRecord, error) {
	if s.returnError == true {
		return nil, fmt.Errorf("fake error")
	}
	return []scraper.FightRecord{}, nil
}

func (d MockDatastore) InsertFightRecords(records []scraper.FightRecord) error {
	if d.returnError == true {
		return fmt.Errorf("fake error")
	}
	return nil
}

func TestHandler_HandleRequest(t *testing.T) {
	tests := []struct {
		name    string
		h       Handler
		wantErr bool
	}{
		{
			name:    "handler should return error due to scraper error",
			h:       Handler{Scraper: MockScraper{returnError: true}, Datastore: MockDatastore{returnError: false}},
			wantErr: true,
		},
		{
			name:    "handler should return error due to datastore error",
			h:       Handler{Scraper: MockScraper{returnError: false}, Datastore: MockDatastore{returnError: true}},
			wantErr: true,
		},
		{
			name:    "handler should not return error",
			h:       Handler{Scraper: MockScraper{returnError: false}, Datastore: MockDatastore{returnError: false}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.HandleRequest(); (err != nil) != tt.wantErr {
				t.Errorf("Handler.HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
