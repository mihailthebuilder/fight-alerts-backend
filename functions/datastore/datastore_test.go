package datastore

import (
	"fight-alerts-backend/scraper"
	utils "fight-alerts-backend/test_utils"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestInsertFightRecords(t *testing.T) {
	d := &Datastore{Host: "localhost", Port: 5432, User: "username", Password: "password", Dbname: "test"}

	err := d.Connect()
	if err != nil {
		t.Errorf("db error - connection: %v", err)
	}

	err = utils.CreateEventTable(d.Db)
	if err != nil {
		t.Errorf("db error - unable to create table: %v", err)
	}

	tz, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Errorf("unable to load London timezone %v", err)
	}

	timeNow := time.Now().Round(time.Second).In(tz)
	timeTomorrow := timeNow.AddDate(0, 0, 1).Round(time.Second).In(tz)
	recordsInsertedInDb := []scraper.FightRecord{{Headline: "today", DateTime: timeNow}, {Headline: "tomorrow", DateTime: timeTomorrow}}

	err = d.InsertFightRecords(recordsInsertedInDb)
	if err != nil {
		t.Errorf("db error - InsertFightRecords throws error: %v", err)
	}

	rows, err := d.Db.Query(`select * from event`)
	if err != nil {
		t.Errorf("db error - select * from event: %v", err)
	}

	var recordsReturnedFromDb []scraper.FightRecord

	for rows.Next() {
		var (
			id       int
			headline string
			date     time.Time
		)

		if err := rows.Scan(&id, &headline, &date); err != nil {
			t.Errorf("db error - scanning through rows returned from query: %v", err)
		}
		recordsReturnedFromDb = append(recordsReturnedFromDb, scraper.FightRecord{DateTime: date, Headline: headline})
	}

	_, err = d.Db.Exec(`drop table Event;`)
	if err != nil {
		t.Errorf("db error - unable drop table: %v", err)
	}

	if !reflect.DeepEqual(recordsInsertedInDb, recordsReturnedFromDb) {
		t.Errorf("records inserted and returned are different. expected: %v . got: %v", recordsInsertedInDb, recordsReturnedFromDb)
	}

	err = d.CloseConnection()
	if err != nil {
		t.Errorf("unable to close db")
	}
}

func TestDatastore_Connect(t *testing.T) {
	tests := []struct {
		name    string
		d       *Datastore
		wantErr bool
	}{
		{
			name:    "valid db connection",
			d:       &Datastore{Host: "localhost", Port: 5432, User: "username", Password: "password", Dbname: "test"},
			wantErr: false,
		},
		{
			name:    "invalid db connection",
			d:       &Datastore{Host: "localhost", Port: 5432, User: "hello", Password: "password", Dbname: "world"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Datastore.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				tt.d.CloseConnection()
			}
		})
	}
}

func TestDatastore_CloseConnection(t *testing.T) {
	d := &Datastore{Host: "localhost", Port: 5432, User: "username", Password: "password", Dbname: "test"}

	err := d.Connect()
	if err != nil {
		t.Errorf("Datastore.CloseConnection() unable to connect")
	}

	err = d.CloseConnection()
	if err != nil {
		t.Errorf("Datastore.CloseConnection() error: %v", err)
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
