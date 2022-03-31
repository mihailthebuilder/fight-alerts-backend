package datastore

import (
	"database/sql"
	"fight-alerts-backend/scraper"
	"fmt"

	"github.com/lib/pq"
)

type Datastore struct {
	db *sql.DB
}

type IDatastore interface {
	InsertFightRecords([]scraper.FightRecord) error
}

func (d Datastore) InsertFightRecords(records []scraper.FightRecord) error {

	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("db error - begin insert: %v", err)
	}

	q, err := tx.Prepare(pq.CopyIn("event", "headline", "datetime"))
	if err != nil {
		return fmt.Errorf("db error - prepare transactions: %v", err)
	}

	for _, record := range records {
		_, err = q.Exec(record.Headline, record.DateTime)
		if err != nil {
			return fmt.Errorf("db error - prepare transactions exec: %v", err)
		}
	}

	_, err = q.Exec()
	if err != nil {
		return fmt.Errorf("db error - transactions exec: %v", err)
	}

	err = q.Close()
	if err != nil {
		return fmt.Errorf("db error - closing transactions: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("db error - commiting transactions: %v", err)
	}

	return nil
}
