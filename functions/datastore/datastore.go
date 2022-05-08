package datastore

import (
	"database/sql"
	"fight-alerts-backend/scraper"
	"fmt"

	"github.com/lib/pq"
)

type Datastore struct {
	Db                           *sql.DB
	Host, Dbname, User, Password string
	Port                         int
}

type IDatastore interface {
	ReplaceWithNewRecords([]scraper.FightRecord) error
}

func (d Datastore) ReplaceWithNewRecords(records []scraper.FightRecord) error {

	err := d.TestConnection()
	if err != nil {
		return err
	}

	err = d.DeleteAllRecords()
	if err != nil {
		return err
	}

	err = d.InsertFightRecords(records)
	if err != nil {
		return err
	}

	err = d.CloseConnection()
	if err != nil {
		return err
	}

	return nil
}

func (d Datastore) InsertFightRecords(records []scraper.FightRecord) error {

	tx, err := d.Db.Begin()
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

func (d Datastore) TestConnection() error {

	err := d.Db.Ping()
	if err != nil {
		return fmt.Errorf("pinging Postgres repository connection: %v", err)
	}

	return nil
}

func (d Datastore) CloseConnection() error {
	return d.Db.Close()
}

func (d Datastore) DeleteAllRecords() error {
	_, err := d.Db.Exec(`DELETE FROM "event"`)

	if err != nil {
		return fmt.Errorf("unable to delete all records: %v", err)
	}

	return nil
}
