package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type AuroraClient struct {
	host string
	port int
	db   *sql.DB
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

func (a *AuroraClient) getAllRowsFromEventTable() ([]interface{}, error) {
	rowCursor, err := a.db.Query(`select * from event`)
	if err != nil {
		return nil, fmt.Errorf("db error - select * from event: %v", err)
	}

	var rows []interface{}

	for rowCursor.Next() {
		var (
			id       int
			headline string
			date     time.Time
		)

		if err := rowCursor.Scan(&id, &headline, &date); err != nil {
			return nil, fmt.Errorf("db error - scanning through rows returned from query: %v", err)
		}
		rows = append(rows, struct {
			date     time.Time
			headline string
		}{date, headline})
	}

	return rows, nil
}
