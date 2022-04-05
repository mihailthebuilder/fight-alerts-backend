package handler

import (
	"fight-alerts-backend/datastore"
	"fight-alerts-backend/scraper"
	"fmt"
	"log"
)

type Handler struct {
	Scraper   scraper.IScraper
	Datastore datastore.IDatastore
}

func (h Handler) HandleRequest() {
	records, err := h.Scraper.GetResultsFromUrl()

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Scraped data:\n%#v\n", records)

	err = h.Datastore.Connect()
	if err != nil {
		log.Panic(err)
	}

	err = h.Datastore.InsertFightRecords(records)
	if err != nil {
		log.Panic(err)
	}

	err = h.Datastore.CloseConnection()
	if err != nil {
		log.Panic(err)
	}
}
