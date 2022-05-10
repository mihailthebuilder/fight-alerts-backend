package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type steps struct {
	containers  Containers
	datastore   AuroraClient
	eventbridge EventBridgeClient
}

var originalFightRecords = []FightRecord{{Headline: "one", DateTime: time.Now()}, {Headline: "two", DateTime: time.Now()}}

func (s *steps) someFightRecordsExistInDb() error {

	err := s.datastore.insertFightRecordsToEventTable(originalFightRecords)
	return err
}

func (s *steps) scraperLambdaIsInvoked() error {
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

func (s *steps) eventBridgeTriggerIsReplaced() error {

	rules, err := s.eventbridge.GetAllRuleNamesByNamePrefix("fight-alerts-notification-rule")
	if err != nil {
		return err
	}

	if len(rules) != 1 {
		return fmt.Errorf("only 1 rule should be created. rules created %#v", rules)
	}

	ruleName := *rules[0].Name

	if ruleName != "fight-alerts-notification-rule" {
		return fmt.Errorf("rule should have correct name, instead it's %v", ruleName)
	}

	targets, err := s.eventbridge.GetAllTargetIdsByRuleName(ruleName)
	if err != nil {
		return err
	}

	if len(targets) != 1 {
		return fmt.Errorf("only 1 target should be created. targets created %#v", targets)
	}

	if *targets[0].Id != "fight-alerts-notification-rule-target-id" ||
		*targets[0].Arn != "arn:aws:lambda:us-east-1:111111111111:function:mock-lambda-arn" {
		return fmt.Errorf("incorrect target settings - id=%v, arn=%v", *targets[0].Id, *targets[0].Arn)

	}

	return nil
}

func (s *steps) triggerForNotificationServiceIsSet() error {

	err := s.eventbridge.InsertRuleWithTarget("rule-to-be-deleted", "target-to-be-deleted", "arn:aws:lambda:us-east-1:9999999999:function:arn-to-be-deleted")
	return err
}
