package handler

import (
	"fight-alerts-backend/scraper"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockScraper struct {
	mock.Mock
}

func (s MockScraper) GetResultsFromUrl() ([]scraper.FightRecord, error) {
	args := s.Called()
	return args.Get(0).([]scraper.FightRecord), args.Error(1)
}

type MockNotificationScheduler struct {
	mock.Mock
}

func (n MockNotificationScheduler) UpdateTrigger(t time.Time) error {
	args := n.Called(t)
	return args.Error(0)
}

type MockDatastore struct {
	mock.Mock
}

func (d MockDatastore) ReplaceWithNewRecords(records []scraper.FightRecord) error {
	args := d.Called(records)
	return args.Error(0)
}
