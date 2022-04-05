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
	insertFightRecordsReturn error
	connectReturn            error
	closeConnectionReturn    error
}

func (s MockScraper) GetResultsFromUrl() ([]scraper.FightRecord, error) {
	if s.returnError == true {
		return nil, fmt.Errorf("fake error")
	}
	return []scraper.FightRecord{}, nil
}

func (d MockDatastore) InsertFightRecords(records []scraper.FightRecord) error {
	return d.insertFightRecordsReturn
}

func (d MockDatastore) Connect() error {
	return d.connectReturn
}

func (d MockDatastore) CloseConnection() error {
	return d.closeConnectionReturn
}

func TestHandler_HandleRequest(t *testing.T) {
	tests := []struct {
		name      string
		h         Handler
		wantPanic bool
	}{
		{
			name: "handler should return error due to scraper error",
			h: Handler{
				Scraper:   MockScraper{returnError: true},
				Datastore: MockDatastore{},
			},
			wantPanic: true,
		},
		{
			name: "handler should return error due to datastore connection error",
			h: Handler{
				Scraper: MockScraper{returnError: false},
				Datastore: MockDatastore{
					connectReturn: fmt.Errorf("Datastore.Connect() fake error"),
				},
			},
			wantPanic: true,
		},
		{
			name: "handler should return error due to datastore insertion error",
			h: Handler{
				Scraper: MockScraper{returnError: false},
				Datastore: MockDatastore{
					insertFightRecordsReturn: fmt.Errorf("Datastore.InsertFightRecords() fake error"),
				},
			},
			wantPanic: true,
		},
		{
			name: "handler should return error due to datastore close error",
			h: Handler{
				Scraper: MockScraper{returnError: false},
				Datastore: MockDatastore{
					closeConnectionReturn: fmt.Errorf("Datastore.Close() fake error"),
				},
			},
			wantPanic: true,
		},
		{
			name:      "handler should not return error",
			h:         Handler{Scraper: MockScraper{returnError: false}, Datastore: MockDatastore{}},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Turn off panic
			defer func() { recover() }()

			tt.h.HandleRequest()

			// Never reaches here if `HandleRequest` panics
			if tt.wantPanic {
				t.Errorf("Handler.HandleRequest() should have panicked")
			}
		})
	}
}
