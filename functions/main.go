package main

import (
	"context"
	"database/sql"
	"fight-alerts-backend/datastore"
	handler "fight-alerts-backend/lambda_handler"
	"fight-alerts-backend/scheduler"
	"fight-alerts-backend/scraper"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

func main() {
	s := scraper.Scraper{Url: scraper.MmaUrl}

	d := setUpDatastore()

	ns := setUpNotificationScheduler()

	lh := handler.Handler{Scraper: s, Datastore: d, NotificationScheduler: ns}

	lambda.Start(lh.HandleRequest)
}

func setUpDatastore() datastore.Datastore {
	host := os.Getenv("RDS_HOST")
	password := os.Getenv("RDS_PASSWORD")
	username := os.Getenv("RDS_USERNAME")

	cs := buildConnectionString(host, password, username, "FightAlertsDb", 5432)
	db, _ := sql.Open("postgres", cs)

	return datastore.Datastore{Db: db}
}

func setUpNotificationScheduler() scheduler.CloudWatchEventsScheduler {
	arn := os.Getenv("NOTIFICATION_LAMBDA_ARN")
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	client := cloudwatchevents.NewFromConfig(cfg)

	return scheduler.CloudWatchEventsScheduler{
		TargetARN:       arn,
		RuleName:        "fight-alerts-notification-rule",
		TargetId:        "fight-alerts-notification-rule-target-id",
		RuleDescription: "Triggers the Notification service for the Fight Alerts product",
		Api:             client,
	}
}

func buildConnectionString(host, password, username, dbName string, port int) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName,
	)
}
