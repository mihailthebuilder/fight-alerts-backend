package main

import (
	"database/sql"
	"fight-alerts-backend/scraper"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type AuroraClient struct {
	host   string
	port   int
	dbconx *sql.DB
}

func (a *AuroraClient) connectToDatabase() *sql.DB {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		a.host, a.port, PostgresConxDetails.User, PostgresConxDetails.Password, PostgresConxDetails.Database,
	)

	dbconx, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	return dbconx
}

func (a *AuroraClient) closeDatabaseConnexion() {
	err := a.dbconx.Close()
	if err != nil {
		panic(err)
	}
}

func (a *AuroraClient) getAllItems() []scraper.FightRecord {

	var err error

	rows, err := a.dbconx.Query(`select * from Event`)
	if err != nil {
		panic(err)
	}

	var records []scraper.FightRecord

	for rows.Next() {
		var (
			id       int
			headline string
			date     time.Time
		)

		if err := rows.Scan(&id, &headline, &date); err != nil {
			panic(err)
		}
		records = append(records, scraper.FightRecord{DateTime: date, Headline: headline})
	}

	return records
}
