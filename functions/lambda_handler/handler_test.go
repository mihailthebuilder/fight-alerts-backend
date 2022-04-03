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
		name    string
		h       Handler
		wantErr bool
	}{
		{
			name: "handler should return error due to scraper error",
			h: Handler{
				Scraper:   MockScraper{returnError: true},
				Datastore: MockDatastore{},
			},
			wantErr: true,
		},
		{
			name: "handler should return error due to datastore connection error",
			h: Handler{
				Scraper: MockScraper{returnError: false},
				Datastore: MockDatastore{
					connectReturn: fmt.Errorf("Datastore.Connect() fake error"),
				},
			},
			wantErr: true,
		},
		{
			name: "handler should return error due to datastore insertion error",
			h: Handler{
				Scraper: MockScraper{returnError: false},
				Datastore: MockDatastore{
					insertFightRecordsReturn: fmt.Errorf("Datastore.InsertFightRecords() fake error"),
				},
			},
			wantErr: true,
		},
		{
			name: "handler should return error due to datastore close error",
			h: Handler{
				Scraper: MockScraper{returnError: false},
				Datastore: MockDatastore{
					closeConnectionReturn: fmt.Errorf("Datastore.Close() fake error"),
				},
			},
			wantErr: true,
		},
		{
			name:    "handler should not return error",
			h:       Handler{Scraper: MockScraper{returnError: false}, Datastore: MockDatastore{}},
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
