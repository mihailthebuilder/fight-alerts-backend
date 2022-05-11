package datastore

import (
	"fight-alerts-backend/scraper-lambda/scraper"
	"fmt"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockedInsertMethodsReturn struct {
	dbBegin, txPrepare, sExec, txCommit error
}

func (mmr MockedInsertMethodsReturn) anyErrors() bool {
	return mmr.dbBegin != nil || mmr.txPrepare != nil || mmr.sExec != nil || mmr.txCommit != nil
}

func setInsertMockExpectations(mock sqlmock.Sqlmock, mmr MockedInsertMethodsReturn, records []scraper.FightRecord) {

	mock.ExpectBegin().WillReturnError(mmr.dbBegin)

	queryInRegex := regexp.QuoteMeta(`COPY "event" ("headline", "datetime") FROM STDIN`)

	mock.ExpectPrepare(queryInRegex).WillReturnError(mmr.txPrepare)

	for _, record := range records {
		if mmr.sExec != nil {
			mock.ExpectExec(queryInRegex).WithArgs(record.Headline, record.DateTime).WillReturnError(mmr.sExec)
		} else {
			mock.ExpectExec(queryInRegex).WithArgs(record.Headline, record.DateTime).WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}

	mock.ExpectCommit().WillReturnError(mmr.txCommit)

	if mmr.anyErrors() {
		mock.ExpectRollback()
	}
}

func createMockFightRecords() ([]scraper.FightRecord, error) {

	tz, err := time.LoadLocation("Europe/London")
	if err != nil {
		return nil, fmt.Errorf("unable to load Europe/London timezone %v", err)
	}

	timeNow := time.Now().Round(time.Second).In(tz)
	timeTomorrow := timeNow.AddDate(0, 0, 1).Round(time.Second).In(tz)
	records := []scraper.FightRecord{{Headline: "today", DateTime: timeNow}, {Headline: "tomorrow", DateTime: timeTomorrow}}

	return records, nil
}
