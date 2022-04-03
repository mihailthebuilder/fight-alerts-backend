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
	Connect() error
	CloseConnection() error
	InsertFightRecords([]scraper.FightRecord) error
}

func (d *Datastore) InsertFightRecords(records []scraper.FightRecord) error {

	tx, err := d.Db.Begin()
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

func (d *Datastore) Connect() error {

	s := d.createDbConnectionString()
	d.Db, _ = sql.Open("postgres", s)

	err := d.Db.Ping()
	if err != nil {
		return fmt.Errorf("pinging Postgres repository connexion at %s: %v", s, err)
	}

	return nil
}

func (d *Datastore) CloseConnection() error {
	return d.Db.Close()
}

func (d *Datastore) createDbConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Dbname,
	)
}
