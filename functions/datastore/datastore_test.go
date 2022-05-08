package datastore

import (
	"fight-alerts-backend/scraper"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatastore_ReplaceWithNewRecords(t *testing.T) {
	tests := []struct {
		name         string
		connectError error
		deleteError  error
		insertError  error
		closeError   error
		wantErr      bool
	}{
		{
			name:    "happy path",
			wantErr: false,
		},
		{
			name:         "connect error",
			connectError: fmt.Errorf("fake err"),
			wantErr:      true,
		},
		{
			name:        "delete error",
			deleteError: fmt.Errorf("fake err"),
			wantErr:     true,
		},
		{
			name:        "insert error",
			insertError: fmt.Errorf("fake err"),
			wantErr:     true,
		},
		{
			name:       "close error",
			closeError: fmt.Errorf("fake err"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			mock.ExpectPing().WillReturnError(tt.connectError)

			q := `DELETE FROM "event"`
			if tt.deleteError != nil {
				mock.ExpectExec(q).WillReturnError(tt.deleteError)
			} else {
				mock.ExpectExec(q).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			records := []scraper.FightRecord{}

			mmr := MockedInsertMethodsReturn{dbBegin: tt.insertError}
			setInsertMockExpectations(mock, mmr, records)

			mock.ExpectClose().WillReturnError(tt.closeError)

			d := Datastore{Db: db}

			if err := d.ReplaceWithNewRecords(records); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.ReplaceWithNewRecords() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil && !tt.wantErr {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatastore_InsertFightRecords(t *testing.T) {

	tests := []struct {
		name    string
		mmr     MockedInsertMethodsReturn
		wantErr bool
	}{
		{name: "happy path", mmr: MockedInsertMethodsReturn{}, wantErr: false},
		{name: "db.Begin error", mmr: MockedInsertMethodsReturn{dbBegin: fmt.Errorf("fake error")}, wantErr: true},
		{name: "tx.Prepare error", mmr: MockedInsertMethodsReturn{txPrepare: fmt.Errorf("fake error")}, wantErr: true},
		{name: "s.Exec error", mmr: MockedInsertMethodsReturn{sExec: fmt.Errorf("fake error")}, wantErr: true},
		{name: "tx.Commit error", mmr: MockedInsertMethodsReturn{txCommit: fmt.Errorf("fake error")}, wantErr: true},
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
			defer db.Close()

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
			defer db.Close()

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

func TestDatastore_DeleteAllRecords(t *testing.T) {
	tests := []struct {
		name               string
		deleteMethodReturn error
		wantErr            bool
	}{
		{name: "happy path", wantErr: false},
		{name: "delete method should return error", deleteMethodReturn: fmt.Errorf("fake error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			q := `DELETE FROM "event"`

			if tt.deleteMethodReturn != nil {
				mock.ExpectExec(q).WillReturnError(tt.deleteMethodReturn)
			} else {
				mock.ExpectExec(q).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			d := &Datastore{Db: db}

			if err := d.DeleteAllRecords(); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.DeleteAllRecords() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
