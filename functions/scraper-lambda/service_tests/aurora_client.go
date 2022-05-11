package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type AuroraClient struct {
	host string
	port int
	db   *sql.DB
}

type FightRecord struct {
	DateTime time.Time
	Headline string
}

func (a *AuroraClient) connectToDatabase() error {

	conxString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		a.host, a.port, PostgresConxDetails.User, PostgresConxDetails.Password, PostgresConxDetails.Database,
	)

	a.db, _ = sql.Open("postgres", conxString)

	err := a.db.Ping()
	if err != nil {
		return fmt.Errorf("pinging Postgres repository connexion at %s: %v", conxString, err)
	}

	return nil
}

func (a *AuroraClient) createEventTable() error {
	_, err := a.db.Exec(`
		create table event (
			event_id serial,
			headline varchar(100),
			datetime timestamptz
		);
	`)
	return err
}

func (a *AuroraClient) getAllFightRecordsFromEventTable() ([]FightRecord, error) {
	rowCursor, err := a.db.Query(`select * from event`)
	if err != nil {
		return nil, fmt.Errorf("db error - select * from event: %v", err)
	}

	var records []FightRecord

	for rowCursor.Next() {
		var (
			id       int
			headline string
			date     time.Time
		)

		if err := rowCursor.Scan(&id, &headline, &date); err != nil {
			return nil, fmt.Errorf("db error - scanning through rows returned from query: %v", err)
		}
		records = append(records, FightRecord{date, headline})
	}

	return records, nil
}

func (a *AuroraClient) insertFightRecordsToEventTable(records []FightRecord) error {

	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("db error - begin insert: %v", err)
	}

	s, err := tx.Prepare(pq.CopyIn("event", "headline", "datetime"))
	if err != nil {
		return fmt.Errorf("db error - prepare transactions: %v", err)
	}

	defer tx.Rollback()

	for _, record := range records {
		_, err = s.Exec(record.Headline, record.DateTime)
		if err != nil {
			return fmt.Errorf("db error - transaction statement exec: %v", err)
		}
	}

	err = s.Close()
	if err != nil {
		return fmt.Errorf("db error - closing transaction statement: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("db error - commiting transactions: %v", err)
	}

	return nil
}
