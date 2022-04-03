package handler

import (
	"fight-alerts-backend/datastore"
	"fight-alerts-backend/scraper"
	"fmt"
)

type Handler struct {
	Scraper   scraper.IScraper
	Datastore datastore.IDatastore
}

func (h Handler) HandleRequest() error {
	records, err := h.Scraper.GetResultsFromUrl()

	if err != nil {
		return err
	}

	fmt.Printf("Scraped data:\n%#v\n", records)

	err = h.Datastore.Connect()
	if err != nil {
		return err
	}

	err = h.Datastore.InsertFightRecords(records)
	if err != nil {
		return err
	}

	err = h.Datastore.CloseConnection()
	if err != nil {
		return err
	}

	return nil
}
