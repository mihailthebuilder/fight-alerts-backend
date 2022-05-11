package handler

import (
	"fight-alerts-backend/scraper-lambda/datastore"
	"fight-alerts-backend/scraper-lambda/scheduler"
	"fight-alerts-backend/scraper-lambda/scraper"
	"fmt"
	"log"
)

type Handler struct {
	Scraper               scraper.IScraper
	Datastore             datastore.IDatastore
	NotificationScheduler scheduler.IScheduler
}

func (h Handler) HandleRequest() {

	records, err := h.Scraper.GetResultsFromUrl()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Scraped data:\n%#v\n", records)

	err = h.NotificationScheduler.UpdateTrigger(records[0].DateTime)
	if err != nil {
		log.Panic(err)
	}

	err = h.Datastore.ReplaceWithNewRecords(records)
	if err != nil {
		log.Panic(err)
	}
}
