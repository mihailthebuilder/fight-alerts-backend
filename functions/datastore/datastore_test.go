package datastore

import (
	"fight-alerts-backend/scraper"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockedMethodsReturn struct {
	dbBegin, txPrepare, sExec, txCommit error
}

func (mmr MockedMethodsReturn) anyErrors() bool {
	return mmr.dbBegin != nil || mmr.txPrepare != nil || mmr.sExec != nil || mmr.txCommit != nil
}

var ErrFake = fmt.Errorf("fake error")

func TestDatastore_InsertFightRecords(t *testing.T) {

	tests := []struct {
		name    string
		mmr     MockedMethodsReturn
		wantErr bool
	}{
		{name: "happy path", mmr: MockedMethodsReturn{}, wantErr: false},
		{name: "db.Begin error", mmr: MockedMethodsReturn{dbBegin: ErrFake}, wantErr: true},
		{name: "tx.Prepare error", mmr: MockedMethodsReturn{txPrepare: ErrFake}, wantErr: true},
		{name: "s.Exec error", mmr: MockedMethodsReturn{sExec: ErrFake}, wantErr: true},
		{name: "tx.Commit error", mmr: MockedMethodsReturn{txCommit: ErrFake}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			records, err := createMockFightRecords()
			if err != nil {
				t.Errorf("Error creating mock fight records: %v", err)
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			setInsertMockExpectations(mock, tt.mmr, records)

			d := &Datastore{Db: db}

			if err := d.InsertFightRecords(records); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.InsertFightRecords() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil && !tt.mmr.anyErrors() {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatastore_Connect_ShouldReturnError(t *testing.T) {

	d := &Datastore{}
	err := d.Connect()

	if err == nil {
		t.Errorf("Datastore.Connect() should return error. Datastore = %v", d)
	}
}

func TestDatastore_CloseConnection(t *testing.T) {

	tests := []struct {
		name                 string
		datastoreCloseReturn error
		wantErr              bool
	}{
		{
			name:                 "valid db close should return nil",
			datastoreCloseReturn: nil,
			wantErr:              false,
		},
		{
			name:                 "invalid db close should return error",
			datastoreCloseReturn: fmt.Errorf("fake error"),
			wantErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			mock.ExpectClose().WillReturnError(tt.datastoreCloseReturn)

			d := &Datastore{Db: db}

			if err := d.CloseConnection(); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatastore_createDbConnectionString(t *testing.T) {
	d := &Datastore{Host: "localhost", Port: 5432, User: "username", Password: "password", Dbname: "test"}
	got := d.createDbConnectionString()

	expected := "host=localhost port=5432 user=username password=password dbname=test sslmode=disable"
	if got != expected {
		t.Errorf("Datastore.createDbConnectionString() = %v, want = %v", got, expected)
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

func setInsertMockExpectations(mock sqlmock.Sqlmock, mmr MockedMethodsReturn, records []scraper.FightRecord) {

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
