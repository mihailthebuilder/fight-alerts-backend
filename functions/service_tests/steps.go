package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type steps struct {
	containers Containers
	datastore  AuroraClient
}

var originalFightRecords = []FightRecord{{Headline: "one", DateTime: time.Now()}, {Headline: "two", DateTime: time.Now()}}

func (s *steps) fightRecordsExistInDb() error {

	err := s.datastore.insertFightRecordsToEventTable(originalFightRecords)
	return err
}

func (s *steps) lambdaIsInvoked() error {
	port, err := s.containers.GetLocalhostPort(s.containers.lambdaContainer, LambdaPort)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d/2015-03-31/functions/myfunction/invocations", s.datastore.host, port)

	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body := buf.String()

	if response.StatusCode != 200 {
		return fmt.Errorf("invoking Lambda: %d %s", response.StatusCode, body)
	}

	return nil
}

func (s *steps) originalFightRecordsAreDeleted() error {

	recordsInDb, err := s.datastore.getAllFightRecordsFromEventTable()

	if err != nil {
		return err
	}

	for _, old := range originalFightRecords {
		for _, new := range recordsInDb {
			if old.Headline == new.Headline {
				return fmt.Errorf("original fight record %v wasn't deleted from the database", old)
			}
		}
	}

	return nil
}

func (s *steps) newFightRecordsAreInserted() error {

	records, err := s.datastore.getAllFightRecordsFromEventTable()

	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no fight records were inserted into the database")
	}

	return nil
}
