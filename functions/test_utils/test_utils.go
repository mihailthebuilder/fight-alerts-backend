package utils_test

import (
	"database/sql"
	"fight-alerts-backend/scraper"
	"fmt"
	"time"
)

func CreateEventTable(db *sql.DB) error {
	_, err := db.Exec(`
		create table event (
			event_id serial,
			headline varchar(100),
			datetime timestamptz
		);
	`)
	return err
}

func GetAllFightRecordsFromEventTable(db *sql.DB) ([]scraper.FightRecord, error) {
	rows, err := db.Query(`select * from event`)
	if err != nil {
		return nil, fmt.Errorf("db error - select * from event: %v", err)
	}

	var records []scraper.FightRecord

	for rows.Next() {
		var (
			id       int
			headline string
			date     time.Time
		)

		if err := rows.Scan(&id, &headline, &date); err != nil {
			return nil, fmt.Errorf("db error - scanning through rows returned from query: %v", err)
		}
		records = append(records, scraper.FightRecord{DateTime: date, Headline: headline})
	}

	return records, nil
}
