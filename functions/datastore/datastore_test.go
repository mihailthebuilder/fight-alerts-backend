package datastore

import (
	"database/sql"
	"fight-alerts-backend/scraper"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestInsertFightRecords(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://username:password@localhost/test?sslmode=disable")

	if err != nil {
		t.Errorf("db error - connection: %v", err)
	}

	defer db.Close()

	_, err = db.Exec(`
		create table event (
			event_id serial,
			headline varchar(100),
			datetime timestamptz
		);
	`)
	if err != nil {
		t.Errorf("db error - unable to create table: %v", err)
	}

	d := Datastore{db: db}

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

	rows, err := d.db.Query(`select * from Event`)
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

	_, err = db.Exec(`drop table Event;`)
	if err != nil {
		t.Errorf("db error - unable drop table: %v", err)
	}

	if !reflect.DeepEqual(recordsInsertedInDb, recordsReturnedFromDb) {
		t.Errorf("records inserted and returned are different. expected: %v . got: %v", recordsInsertedInDb, recordsReturnedFromDb)
	}
}
