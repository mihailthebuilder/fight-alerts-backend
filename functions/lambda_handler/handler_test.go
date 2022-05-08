package handler

import (
	"fight-alerts-backend/scraper"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler_HandleRequest(t *testing.T) {
	tests := []struct {
		name           string
		scraperError   error
		datastoreError error
		schedulerError error
		wantErr        bool
	}{
		{
			name:    "happy path",
			wantErr: false,
		},
		{
			name:         "scraper error",
			scraperError: fmt.Errorf("fake err"),
			wantErr:      true,
		},
		{
			name:           "scraper error",
			datastoreError: fmt.Errorf("fake err"),
			wantErr:        true,
		},
		{
			name:           "scraper error",
			schedulerError: fmt.Errorf("fake err"),
			wantErr:        true,
		},
	}
	for _, tt := range tests {

		firstRecord := scraper.FightRecord{Headline: "first", DateTime: time.Now()}
		secondRecord := scraper.FightRecord{Headline: "second", DateTime: time.Now()}
		records := []scraper.FightRecord{firstRecord, secondRecord}

		s := MockScraper{}
		s.On("GetResultsFromUrl").Return(records, tt.scraperError)

		n := MockNotificationScheduler{}
		n.On("UpdateTrigger", firstRecord.DateTime).Return(tt.schedulerError)

		d := MockDatastore{}
		d.On("ReplaceWithNewRecords", records).Return(tt.datastoreError)

		h := Handler{Scraper: s, Datastore: d, NotificationScheduler: n}

		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, h.HandleRequest, "The code did not panic")
			} else {
				h.HandleRequest()
				s.AssertExpectations(t)
				n.AssertExpectations(t)
				d.AssertExpectations(t)
			}
		})
	}
}
